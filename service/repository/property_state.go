package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
)

type PropertyStateRepository interface {
	GetByID(id string) (*models.PropertyState, error)
	GetByPropertyID(id string) (*models.PropertyState, error)
	GetAllByPropertyID(id string) ([]models.PropertyState, error)
	Save(propertyState *models.PropertyState) error
}

type propertyStateRepository struct {
	readDb  *gorm.DB
	writeDb *gorm.DB
}

func NewPropertyStateRepository(ctx context.Context, service *frame.Service) PropertyStateRepository {
	return &propertyStateRepository{readDb: service.DB(ctx, true), writeDb: service.DB(ctx, false)}
}

func (repo *propertyStateRepository) GetByPropertyID(id string) (*models.PropertyState, error) {
	var propertyState models.PropertyState
	err := repo.readDb.Last(&propertyState, "property_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &propertyState, nil
}

func (repo *propertyStateRepository) GetAllByPropertyID(id string) ([]models.PropertyState, error) {
	var propertyStates []models.PropertyState
	err := repo.readDb.Find(&propertyStates, "property_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return propertyStates, nil
}

func (repo *propertyStateRepository) GetByID(id string) (*models.PropertyState, error) {
	propertyState := models.PropertyState{}
	err := repo.readDb.First(&propertyState, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &propertyState, nil
}

func (repo *propertyStateRepository) Save(propertyState *models.PropertyState) error {
	err := repo.writeDb.Save(propertyState).Error
	if frame.DBErrorIsRecordNotFound(err) {
		return repo.writeDb.Create(propertyState).Error
	}
	return nil
}
