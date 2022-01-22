package repository

import (
	"fmt"
	"github.com/pitabwire/frame"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type baseRepository struct {
	readDb          *gorm.DB
	writeDb         *gorm.DB
	instanceCreator func() frame.BaseModelI
}

func (repo *baseRepository) getReadDb() *gorm.DB {
	return repo.readDb
}

func (repo *baseRepository) getWriteDb() *gorm.DB {
	return repo.writeDb
}

func (repo *baseRepository) Delete(id string) error {
	deleteInstance, err := repo.GetByID(id)
	if err != nil {
		return err
	}

	return repo.writeDb.Delete(deleteInstance).Error

}

func (repo *baseRepository) GetByID(id string) (frame.BaseModelI, error) {
	getInstance := repo.instanceCreator()
	err := repo.readDb.Preload(clause.Associations).First(getInstance, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return getInstance, nil
}

func (repo *baseRepository) GetLastestBy(properties map[string]interface{}) (frame.BaseModelI, error) {
	getInstance := repo.instanceCreator()

	db := repo.readDb

	for key, value := range properties {
		db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	err := db.Last(getInstance).Error
	if err != nil {
		return nil, err
	}
	return getInstance, nil
}

func (repo *baseRepository) GetAllBy(properties map[string]interface{}, instanceList interface{}) error {

	db := repo.readDb

	for key, value := range properties {
		db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	return db.Find(instanceList).Error
}

func (repo *baseRepository) Search(query string, searchFields []string, instanceList interface{}) error {

	db := repo.readDb

	for i, field := range searchFields {
		if i == 0 {
			db.Where(fmt.Sprintf("%s iLike ?", field), query)
		} else {
			db.Or(fmt.Sprintf(" %s iLike ?", field), query)
		}
	}

	return db.Find(instanceList).Error
}

func (repo *baseRepository) Save(instance frame.BaseModelI) error {

	if instance.GetVersion() <= 0 {

		err := repo.writeDb.Create(instance).Error
		if err != nil {
			return err
		}
	} else {
		return repo.writeDb.Save(instance).Error
	}
	return nil
}
