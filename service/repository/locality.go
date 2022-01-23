package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
)

type LocalityRepository interface {
	frame.BaseRepositoryI
}

type localityRepository struct {
	*frame.BaseRepository
}

func NewLocalityRepository(ctx context.Context, service *frame.Service) LocalityRepository {
	return &localityRepository{
		BaseRepository: frame.NewBaseRepository(
			service.DB(ctx, true),
			service.DB(ctx, false),
			func() frame.BaseModelI {
				return &models.Locality{}
			}),
	}
}
