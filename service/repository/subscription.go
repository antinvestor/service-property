package repository

import (
	"context"

	"github.com/antinvestor/service-property/service/models"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	GetByID(id string) (*models.Subscription, error)
	GetByPropertyID(propertyId string) ([]models.Subscription, error)
	Save(subscription *models.Subscription) error
	Delete(id string) error
}

type subscriptionRepository struct {
	readDb  *gorm.DB
	writeDb *gorm.DB
}

func NewSubscriptionRepository(ctx context.Context, service *frame.Service) SubscriptionRepository {
	return &subscriptionRepository{readDb: service.DB(ctx, true), writeDb: service.DB(ctx, false)}
}

func (repo *subscriptionRepository) GetByID(id string) (*models.Subscription, error) {
	subscription := models.Subscription{}
	err := repo.readDb.First(&subscription, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (repo *subscriptionRepository) GetByPropertyID(propertyId string) ([]models.Subscription, error) {
	var subscriptionList []models.Subscription

	err := repo.readDb.Find(&subscriptionList,
		"property_id = ? ", propertyId).Error
	if err != nil {
		return nil, err
	}
	return subscriptionList, nil
}

func (repo *subscriptionRepository) Save(subscription *models.Subscription) error {
	err := repo.writeDb.Save(subscription).Error
	if frame.DBErrorIsRecordNotFound(err) {
		return repo.writeDb.Create(subscription).Error
	}
	return nil
}

func (repo *subscriptionRepository) Delete(id string) error {
	subscription, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	return repo.writeDb.Delete(subscription).Error

}
