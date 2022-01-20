package business

import (
	"context"
	partapi "github.com/antinvestor/service-partition-api"
	propertyV1 "github.com/antinvestor/service-property-api"

	profileV1 "github.com/antinvestor/service-profile-api"
	"github.com/pitabwire/frame"
)

type PropertyBusiness interface {
	CreateProperty(context.Context, *propertyV1.Property) (*propertyV1.PropertyState, error)
	UpdateProperty(context.Context, *propertyV1.UpdateRequest) (*propertyV1.Property, error)
	DeleteProperty(context.Context, *propertyV1.RequestID) (*propertyV1.PropertyState, error)
	StateOfProperty(context.Context, *propertyV1.RequestID) (*propertyV1.PropertyState, error)
	HistoryOfProperty(*propertyV1.RequestID, propertyV1.PropertyService_HistoryOfPropertyServer) error
	SearchProperty(*propertyV1.SearchRequest, propertyV1.PropertyService_SearchPropertyServer) error
}

func NewPropertyBusiness(ctx context.Context, service *frame.Service, profileCli *profileV1.ProfileClient, partitionCli *partapi.PartitionClient) (PropertyBusiness, error) {

	if service == nil || profileCli == nil || partitionCli == nil {
		return nil, ErrorInitializationFail
	}

	return &propertyBusiness{
		service:    service,
		profileCli: profileCli,
	}, nil
}

type PropertyTypeBusiness interface {
	AddPropertyType(context.Context, *propertyV1.PropertyType) (*propertyV1.PropertyType, error)
	ListPropertyType(*propertyV1.SearchRequest, propertyV1.PropertyService_ListTypeServer) error
}

func NewPropertyTypeBusiness(ctx context.Context, service *frame.Service, profileCli *profileV1.ProfileClient, partitionCli *partapi.PartitionClient) (PropertyTypeBusiness, error) {

	if service == nil || profileCli == nil || partitionCli == nil {
		return nil, ErrorInitializationFail
	}

	return &propertyTypeBusiness{
		service:    service,
		profileCli: profileCli,
	}, nil
}

type LocalityBusiness interface {
	AddLocality(context.Context, *propertyV1.Locality) (*propertyV1.Locality, error)
	DeleteLocality(context.Context, *propertyV1.RequestID) error
}

func NewLocalityBusiness(ctx context.Context, service *frame.Service, profileCli *profileV1.ProfileClient, partitionCli *partapi.PartitionClient) (LocalityBusiness, error) {

	if service == nil || profileCli == nil || partitionCli == nil {
		return nil, ErrorInitializationFail
	}

	return &localityBusiness{
		service:    service,
		profileCli: profileCli,
	}, nil
}

type SubscriptionBusiness interface {
	ListSubscriptions(*propertyV1.SubscriptionListRequest, propertyV1.PropertyService_ListSubscriptionsServer) error
	AddSubscription(context.Context, *propertyV1.Subscription) (*propertyV1.Subscription, error)
	DeleteSubscription(context.Context, *propertyV1.RequestID) error
}

func NewSubscriptionBusiness(ctx context.Context, service *frame.Service, profileCli *profileV1.ProfileClient, partitionCli *partapi.PartitionClient) (SubscriptionBusiness, error) {

	if service == nil || profileCli == nil || partitionCli == nil {
		return nil, ErrorInitializationFail
	}

	return &subscriptionBusiness{
		service:    service,
		profileCli: profileCli,
	}, nil
}
