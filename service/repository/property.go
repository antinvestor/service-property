package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
)

type PropertyRepository interface {
	GetByID(id string) (*models.Property, error)
	SearchByPartition(partitionId string, query string) ([]models.Property, error)
	Save(property *models.Property) error
	Delete(id string) error
}

type propertyRepository struct {
	readDb  *gorm.DB
	writeDb *gorm.DB
}

func NewPropertyRepository(ctx context.Context, service *frame.Service) PropertyRepository {
	return &propertyRepository{readDb: service.DB(ctx, true), writeDb: service.DB(ctx, false)}
}

func (repo *propertyRepository) GetByID(id string) (*models.Property, error) {
	property := models.Property{}
	err := repo.readDb.First(&property, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &property, nil
}

func (repo *propertyRepository) SearchByPartition(partitionId string, query string) ([]models.Property, error) {
	var properties []models.Property

	err := repo.readDb.Find(&properties,
		"partition_id = ? AND (id ILIKE ? OR name ILIKE ? OR description ILIKE ?)",
		partitionId, query, query, query).Error
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (repo *propertyRepository) Delete(id string) error {
	property, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	return repo.writeDb.Delete(property).Error

}

func (repo *propertyRepository) Save(property *models.Property) error {
	err := repo.writeDb.Save(property).Error
	if frame.DBErrorIsRecordNotFound(err) {
		return repo.writeDb.Create(property).Error
	}
	return nil
}
