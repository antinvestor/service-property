package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
)

type LocalityRepository interface {
	GetByID(id string) (*models.Locality, error)
	Delete(id string) error
	Save(locality frame.BaseModelI) error
}

type localityRepository struct {
	baseRepository
}

func NewLocalityRepository(ctx context.Context, service *frame.Service) LocalityRepository {
	return &localityRepository{
		baseRepository: baseRepository{
			readDb:  service.DB(ctx, true),
			writeDb: service.DB(ctx, false),
			instanceCreator: func() frame.BaseModelI {
				return &models.Locality{}
			},
		},
	}
}

func (repo *localityRepository) Delete(id string) error {
	locality, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	return repo.writeDb.Delete(locality).Error

}

func (repo *localityRepository) GetByID(id string) (*models.Locality, error) {
	locality, err := repo.baseRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return locality.(*models.Locality), err
}
