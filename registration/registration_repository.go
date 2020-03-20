package registration

import (
	"github.com/jinzhu/gorm"
	"kpdigisign/model"
)

type Repository interface {
	FindAll() ([]model.Registration, error)
	FindById(string) (model.Registration, error)
	Save(model.Registration) (model.Registration, error)
	Update(model.Registration) (model.Registration, error)
	Delete(model.Registration) error
}

type registrationRepository struct {
	*gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) Repository {
	return &registrationRepository{DB: db}
}

func (r registrationRepository) FindAll() ([]model.Registration, error) {
	var entities []model.Registration
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r registrationRepository) FindById(id string) (model.Registration, error) {
	var entity model.Registration
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r registrationRepository) Save(entity model.Registration) (model.Registration, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r registrationRepository) Update(entity model.Registration) (model.Registration, error) {
	err := r.DB.Model(model.Registration{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r registrationRepository) Delete(entity model.Registration) error {
	return r.DB.Delete(&entity).Error
}
