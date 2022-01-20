package business

import (
	"context"
	profileV1 "github.com/antinvestor/service-profile-api"
	propertyV1 "github.com/antinvestor/service-property-api"
	"github.com/antinvestor/service-property/service/models"
	"github.com/antinvestor/service-property/service/repository"
	"github.com/pitabwire/frame"
)

type propertyTypeBusiness struct {
	service    *frame.Service
	profileCli *profileV1.ProfileClient
}

func (pt *propertyTypeBusiness) AddPropertyType(ctx context.Context, message *propertyV1.PropertyType) (*propertyV1.PropertyType, error) {

	err := message.Validate()
	if err != nil {
		return nil, err
	}

	propertyTypeRepo := repository.NewPropertyTypeRepository(ctx, pt.service)

	propertyType := models.PropertyType{
		Name:        message.GetName(),
		Description: message.GetDescription(),
		Extra:       frame.DBPropertiesFromMap(message.GetExtra()),
	}

	if propertyType.ValidXID(message.GetID()) {
		propertyType.ID = message.GetID()
	} else {
		propertyType.GenID(ctx)
	}

	err = propertyTypeRepo.Save(&propertyType)
	if err != nil {
		return nil, err
	}

	return propertyType.ToApi(), nil
}

func (pt *propertyTypeBusiness) ListPropertyType(message *propertyV1.SearchRequest, stream propertyV1.PropertyService_ListTypeServer) error {

	err := message.Validate()
	if err != nil {
		return err
	}

	propertyTypeRepository := repository.NewPropertyTypeRepository(stream.Context(), pt.service)

	propertyTypeList, err := propertyTypeRepository.GetAllByQuery(message.GetQuery())
	if err != nil {
		return err
	}

	for _, propertyType := range propertyTypeList {
		err := stream.Send(propertyType.ToApi())
		if err != nil {
			pt.service.L().Info(" ListPropertyType -- unable to send a result see %v", err)
		}
	}

	return nil
}
