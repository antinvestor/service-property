package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/antinvestor/service-property/service/models"
	"github.com/antinvestor/service-property/service/repository"
	"github.com/pitabwire/frame"
)

type subscriptionBusiness struct {
	service    *frame.Service
	profileCli *profileV1.ProfileClient
}

func (s *subscriptionBusiness) AddSubscription(ctx context.Context, message *propertyV1.Subscription) (*propertyV1.Subscription, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	subscriptionRepository := repository.NewSubscriptionRepository(ctx, s.service)

	locality := models.Subscription{

		PropertyID: message.GetPropertyID(),
		ProfileID:  message.GetProfileID(),
		Role:       message.GetRole(),
		Extra:      frame.DBPropertiesFromMap(message.GetExtra()),
		ExpiresAt:  message.GetExpiresAt().AsTime(),
	}

	if locality.ValidXID(message.GetID()) {
		locality.ID = message.GetID()
	} else {
		locality.GenID(ctx)
	}

	err = subscriptionRepository.Save(&locality)
	if err != nil {
		return nil, err
	}

	return locality.ToApi(), nil
}

func (s *subscriptionBusiness) ListSubscriptions(message *propertyV1.SubscriptionListRequest, stream propertyV1.PropertyService_ListSubscriptionsServer) error {

	err := message.Validate()
	if err != nil {
		return err
	}

	subscriptionRepository := repository.NewSubscriptionRepository(stream.Context(), s.service)

	subscriptionList, err := subscriptionRepository.GetByPropertyID(message.GetPropertyID())
	if err != nil {
		return err
	}

	for _, subscription := range subscriptionList {
		err := stream.Send(subscription.ToApi())
		if err != nil {
			s.service.L().Info(" ListSubscriptions -- unable to send a result see %v", err)
		}
	}

	return nil
}

func (s *subscriptionBusiness) DeleteSubscription(ctx context.Context, message *propertyV1.RequestID) error {

	err := message.Validate()
	if err != nil {
		return err
	}

	subscriptionRepository := repository.NewSubscriptionRepository(ctx, s.service)

	subscription, err := subscriptionRepository.GetByID(message.GetID())
	if err != nil {
		return err
	}

	return subscriptionRepository.Delete(subscription.GetID())
}
