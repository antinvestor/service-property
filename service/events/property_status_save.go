package events

import (
	"context"
	"errors"
	"github.com/antinvestor/service-property/service/models"
	"github.com/antinvestor/service-property/service/repository"
	"github.com/pitabwire/frame"
)

type PropertyStateSave struct {
	Service *frame.Service
}

func (e *PropertyStateSave) Name() string {
	return "propertyState.save"
}

func (e *PropertyStateSave) PayloadType() interface{} {
	return &models.PropertyState{}
}

func (e *PropertyStateSave) Validate(_ context.Context, payload interface{}) error {
	propertyState, ok := payload.(*models.PropertyState)
	if !ok {
		return errors.New(" payload is not of type models.PropertyState")
	}

	if propertyState.GetID() == "" {
		return errors.New(" propertyState Id should already have been set ")
	}

	return nil
}

func (e *PropertyStateSave) Execute(ctx context.Context, payload interface{}) error {

	propertyState := payload.(*models.PropertyState)

	err := e.Service.DB(ctx, false).Save(propertyState).Error

	if err != nil {
		return err
	}

	propertyRepository := repository.NewPropertyRepository(ctx, e.Service)
	property, err := propertyRepository.GetByID(propertyState.PropertyID)
	if err != nil {
		return err
	}

	property.StateID = propertyState.ID

	err = propertyRepository.Save(property)
	if err != nil {
		return err
	}

	return nil
}
