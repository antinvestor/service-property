package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
)

type PropertyTypeRepository interface {
	GetByID(id string) (*models.PropertyType, error)
	GetAllByQuery(partitionId string) ([]models.PropertyType, error)
	Save(propertyType *models.PropertyType) error
}

type propertyTypeRepository struct {
	readDb  *gorm.DB
	writeDb *gorm.DB
}

func NewPropertyTypeRepository(ctx context.Context, service *frame.Service) PropertyTypeRepository {
	return &propertyTypeRepository{readDb: service.DB(ctx, true), writeDb: service.DB(ctx, false)}
}

func (repo *propertyTypeRepository) GetByID(id string) (*models.PropertyType, error) {
	propertyType := models.PropertyType{}
	err := repo.readDb.First(&propertyType, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &propertyType, nil
}

func (repo *propertyTypeRepository) GetAllByQuery(query string) ([]models.PropertyType, error) {
	var routes []models.PropertyType

	if query == "" {
		err := repo.readDb.Find(&routes).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := repo.readDb.Find(&routes, "name iLike ?", query).Error
		if err != nil {
			return nil, err
		}
	}
	return routes, nil
}

func (repo *propertyTypeRepository) Save(propertyType *models.PropertyType) error {
	err := repo.writeDb.Save(propertyType).Error
	if frame.DBErrorIsRecordNotFound(err) {
		return repo.writeDb.Create(propertyType).Error
	}
	return nil
}
