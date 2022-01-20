package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LocalityRepository interface {
	GetByID(id string) (*models.Locality, error)
	Delete(id string) error
	Save(locality *models.Locality) error
}

type localityRepository struct {
	readDb  *gorm.DB
	writeDb *gorm.DB
}

func NewLocalityRepository(ctx context.Context, service *frame.Service) LocalityRepository {
	return &localityRepository{readDb: service.DB(ctx, true), writeDb: service.DB(ctx, false)}
}

func (repo *localityRepository) Delete(id string) error {
	locality, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	return repo.writeDb.Delete(locality).Error

}

func (repo *localityRepository) GetByID(id string) (*models.Locality, error) {
	locality := models.Locality{}
	err := repo.readDb.Preload(clause.Associations).First(&locality, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &locality, nil
}

func (repo *localityRepository) Save(template *models.Locality) error {
	err := repo.writeDb.Save(template).Error
	if frame.DBErrorIsRecordNotFound(err) {
		return repo.writeDb.Create(template).Error
	}
	return nil
}
