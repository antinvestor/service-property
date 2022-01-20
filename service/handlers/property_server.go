package handlers

import (
	"context"
	partapi "github.com/antinvestor/service-partition-api"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/antinvestor/service-property/service/business"
	"github.com/pitabwire/frame"
)

type PropertyServer struct {
	Service      *frame.Service
	ProfileCli   *profileV1.ProfileClient
	PartitionCli *partapi.PartitionClient

	propertyV1.UnimplementedPropertyServiceServer
}

func (server *PropertyServer) newPropertyBusiness(ctx context.Context) (business.PropertyBusiness, error) {
	return business.NewPropertyBusiness(ctx, server.Service, server.ProfileCli, server.PartitionCli)
}

func (server *PropertyServer) newPropertyTypeBusiness(ctx context.Context) (business.PropertyTypeBusiness, error) {
	return business.NewPropertyTypeBusiness(ctx, server.Service, server.ProfileCli, server.PartitionCli)
}

func (server *PropertyServer) newLocalityBusiness(ctx context.Context) (business.LocalityBusiness, error) {
	return business.NewLocalityBusiness(ctx, server.Service, server.ProfileCli, server.PartitionCli)
}

func (server *PropertyServer) newSubscriptionBusiness(ctx context.Context) (business.SubscriptionBusiness, error) {
	return business.NewSubscriptionBusiness(ctx, server.Service, server.ProfileCli, server.PartitionCli)
}

func (server *PropertyServer) AddPropertyType(ctx context.Context, message *propertyV1.PropertyType) (*propertyV1.PropertyType, error) {

	propertyTypeBusiness, err := server.newPropertyTypeBusiness(ctx)
	if err != nil {
		return nil, err
	}

	return propertyTypeBusiness.AddPropertyType(ctx, message)
}
func (server *PropertyServer) ListType(message *propertyV1.SearchRequest, stream propertyV1.PropertyService_ListTypeServer) error {
	propertyTypeBusiness, err := server.newPropertyTypeBusiness(stream.Context())
	if err != nil {
		return err
	}

	return propertyTypeBusiness.ListPropertyType(message, stream)
}
func (server *PropertyServer) AddLocality(ctx context.Context, message *propertyV1.Locality) (*propertyV1.Locality, error) {
	localityBusiness, err := server.newLocalityBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return localityBusiness.AddLocality(ctx, message)
}
func (server *PropertyServer) DeleteLocality(ctx context.Context, message *propertyV1.RequestID) (*propertyV1.Locality, error) {
	localityBusiness, err := server.newLocalityBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return nil, localityBusiness.DeleteLocality(ctx, message)
}
func (server *PropertyServer) CreateProperty(ctx context.Context, message *propertyV1.Property) (*propertyV1.PropertyState, error) {
	propertyBusiness, err := server.newPropertyBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return propertyBusiness.CreateProperty(ctx, message)
}
func (server *PropertyServer) UpdateProperty(ctx context.Context, message *propertyV1.UpdateRequest) (*propertyV1.Property, error) {
	propertyBusiness, err := server.newPropertyBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return propertyBusiness.UpdateProperty(ctx, message)
}
func (server *PropertyServer) DeleteProperty(ctx context.Context, message *propertyV1.RequestID) (*propertyV1.PropertyState, error) {
	propertyBusiness, err := server.newPropertyBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return propertyBusiness.DeleteProperty(ctx, message)
}
func (server *PropertyServer) StateOfProperty(ctx context.Context, message *propertyV1.RequestID) (*propertyV1.PropertyState, error) {
	propertyBusiness, err := server.newPropertyBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return propertyBusiness.StateOfProperty(ctx, message)
}
func (server *PropertyServer) HistoryOfProperty(message *propertyV1.RequestID, stream propertyV1.PropertyService_HistoryOfPropertyServer) error {
	propertyBusiness, err := server.newPropertyBusiness(stream.Context())
	if err != nil {
		return err
	}
	return propertyBusiness.HistoryOfProperty(message, stream)
}
func (server *PropertyServer) SearchProperty(message *propertyV1.SearchRequest, stream propertyV1.PropertyService_SearchPropertyServer) error {
	propertyBusiness, err := server.newPropertyBusiness(stream.Context())
	if err != nil {
		return err
	}
	return propertyBusiness.SearchProperty(message, stream)
}
func (server *PropertyServer) ListSubscriptions(message *propertyV1.SubscriptionListRequest, stream propertyV1.PropertyService_ListSubscriptionsServer) error {
	subscriptionBusiness, err := server.newSubscriptionBusiness(stream.Context())
	if err != nil {
		return err
	}

	return subscriptionBusiness.ListSubscriptions(message, stream)
}
func (server *PropertyServer) AddSubscription(ctx context.Context, message *propertyV1.Subscription) (*propertyV1.Subscription, error) {
	subscriptionBusiness, err := server.newSubscriptionBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return subscriptionBusiness.AddSubscription(ctx, message)
}
func (server *PropertyServer) DeleteSubscription(ctx context.Context, message *propertyV1.RequestID) (*propertyV1.PropertyState, error) {
	subscriptionBusiness, err := server.newSubscriptionBusiness(ctx)
	if err != nil {
		return nil, err
	}
	return nil, subscriptionBusiness.DeleteSubscription(ctx, message)
}
