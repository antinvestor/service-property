package business

import (
	"context"
	"github.com/antinvestor/apis/common"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/antinvestor/service-property/service/events"
	"github.com/antinvestor/service-property/service/models"
	"github.com/antinvestor/service-property/service/repository"
	"github.com/pitabwire/frame"
	"google.golang.org/protobuf/types/known/timestamppb"
	"runtime"
)

type propertyBusiness struct {
	service    *frame.Service
	profileCli *profileV1.ProfileClient
}

func (pb *propertyBusiness) ToApi(ctx context.Context, property *models.Property) (*propertyV1.Property, error) {

	apiProperty := propertyV1.Property{
		ID:          property.GetID(),
		ParentID:    property.ParentID,
		Name:        property.Name,
		Description: property.Description,
		Extra:       frame.DBPropertiesToMap(property.Extra),
		StartedAt:   timestamppb.New(property.StartedAt),
	}

	if property.LocalityID != "" {
		localityRepo := repository.NewLocalityRepository(ctx, pb.service)

		var locality models.Locality
		err := localityRepo.GetByID(property.LocalityID, &locality)
		if err != nil {
			return nil, err
		}
		apiProperty.Locality = locality.ToApi()
	}

	if property.PropertyTypeID != "" {
		pTypeRepo := repository.NewPropertyTypeRepository(ctx, pb.service)
		propertyType, err := pTypeRepo.GetByID(property.PropertyTypeID)
		if err != nil {
			return nil, err
		}

		apiProperty.PropertyType = propertyType.ToApi()
	}

	return &apiProperty, nil
}

func (pb *propertyBusiness) CreateProperty(ctx context.Context, message *propertyV1.Property) (*propertyV1.PropertyState, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	propertyRepo := repository.NewPropertyRepository(ctx, pb.service)

	property := models.Property{

		Name:           message.GetName(),
		ParentID:       message.GetParentID(),
		PropertyTypeID: message.GetPropertyType().GetID(),
		Description:    message.GetDescription(),
		Extra:          frame.DBPropertiesFromMap(message.GetExtra()),
		StartedAt:      message.GetStartedAt().AsTime(),
	}

	if property.ValidXID(message.GetID()) {
		property.ID = message.GetID()
	} else {
		property.GenID(ctx)
	}

	err = propertyRepo.Save(&property)
	if err != nil {
		return nil, err
	}

	propertyStateRepo := repository.NewPropertyStateRepository(ctx, pb.service)

	propertyState := models.PropertyState{
		PropertyID: property.GetID(),
		State:      int32(common.STATE_CREATED.Number()),
		Status:     int32(common.STATUS_QUEUED.Number()),
		Name:       common.STATE_CREATED.String(),
	}

	propertyState.GenID(ctx)

	err = propertyStateRepo.Save(&propertyState)
	if err != nil {
		return nil, err
	}
	//// Queue out notification status for further processing
	//eventState := events.PropertyStateSave{}
	//err = pb.service.Emit(ctx, eventState.Name(), propertyState)
	//if err != nil {
	//	return nil, err
	//}

	return propertyState.ToApi(), nil
}

func (pb *propertyBusiness) UpdateProperty(ctx context.Context, message *propertyV1.UpdateRequest) (*propertyV1.Property, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	propertyRepo := repository.NewPropertyRepository(ctx, pb.service)
	property, err := propertyRepo.GetByID(message.GetID())
	if err != nil {
		return nil, err
	}

	if message.GetName() != "" {
		property.Name = message.GetName()
	}

	if message.GetDescription() != "" {
		property.Description = message.GetDescription()
	}

	if message.GetExtras() != nil {
		extras := frame.DBPropertiesToMap(property.Extra)
		for key, val := range message.GetExtras() {

			extras[key] = val

		}
		property.Extra = frame.DBPropertiesFromMap(extras)
	}

	err = propertyRepo.Save(property)
	if err != nil {
		return nil, err
	}

	return pb.ToApi(ctx, property)
}

func (pb *propertyBusiness) DeleteProperty(ctx context.Context, message *propertyV1.RequestID) (*propertyV1.PropertyState, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	propertyRepository := repository.NewPropertyRepository(ctx, pb.service)
	property, err := propertyRepository.GetByID(message.GetID())
	if err != nil {
		return nil, err
	}

	err = propertyRepository.Delete(property.ID)
	if err != nil {
		return nil, err
	}

	propertyState := models.PropertyState{
		PropertyID: property.GetID(),
		State:      int32(common.STATE_CREATED.Number()),
		Status:     int32(common.STATUS_QUEUED.Number()),
		Name:       common.STATE_CREATED.String(),
	}

	propertyState.GenID(ctx)

	// Queue property state for further processing
	eventState := events.PropertyStateSave{}
	err = pb.service.Emit(ctx, eventState.Name(), propertyState)
	if err != nil {
		return nil, err
	}

	return propertyState.ToApi(), nil
}

func (pb *propertyBusiness) StateOfProperty(ctx context.Context, message *propertyV1.RequestID) (*propertyV1.PropertyState, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	propertyRepository := repository.NewPropertyRepository(ctx, pb.service)
	property, err := propertyRepository.GetByID(message.GetID())
	if err != nil {
		return nil, err
	}

	propertyStateRepository := repository.NewPropertyStateRepository(ctx, pb.service)

	propertyState, err := propertyStateRepository.GetByPropertyID(property.GetID())
	if err != nil {
		return nil, err
	}

	return propertyState.ToApi(), nil
}

func (pb *propertyBusiness) HistoryOfProperty(message *propertyV1.RequestID, stream propertyV1.PropertyService_HistoryOfPropertyServer) error {

	err := message.Validate()
	if err != nil {
		return err
	}

	propertyRepository := repository.NewPropertyRepository(stream.Context(), pb.service)
	property, err := propertyRepository.GetByID(message.GetID())
	if err != nil {
		return err
	}

	propertyStateRepository := repository.NewPropertyStateRepository(stream.Context(), pb.service)

	propertyStateList, err := propertyStateRepository.GetAllByPropertyID(property.GetID())
	if err != nil {
		return err
	}

	for _, propertyState := range propertyStateList {
		err := stream.Send(propertyState.ToApi())
		if err != nil {
			pb.service.L().Info(" HistoryOfProperty -- unable to send a result see %v", err)
		}
	}

	return nil
}

func (pb *propertyBusiness) SearchProperty(search *propertyV1.SearchRequest, stream propertyV1.PropertyService_SearchPropertyServer) error {

	err := search.Validate()
	if err != nil {
		return err
	}

	propertyRepository := repository.NewPropertyRepository(stream.Context(), pb.service)

	propertyList, err := propertyRepository.Search(search.GetQuery())
	if err != nil {
		return err
	}

	for _, property := range propertyList {

		apiProperty, err := pb.ToApi(stream.Context(), &property)
		if err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			pb.service.L().Info(" SearchProperty -- unable to convert a result : %s", buf)
		}
		err = stream.Send(apiProperty)
		if err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			pb.service.L().Info(" SearchProperty -- unable to send a result : %s", buf)
		}
	}

	return nil
}
