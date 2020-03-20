package activation

import (
	"github.com/jinzhu/gorm"
	"kpdigisign/model"
)

type Repository interface {
	FindAll() ([]model.Activation, error)
	FindById(string) (model.Activation, error)
	Save(model.Activation) (model.Activation, error)
	SaveCallback(callback model.ActivationCallback) (model.ActivationCallback, error)
	Update(model.Activation) (model.Activation, error)
	Delete(model.Activation) error
}

type activationRepository struct {
	*gorm.DB
}

func NewActivationRepository(db *gorm.DB) Repository {
	return &activationRepository{DB: db}
}

func (r activationRepository) FindAll() ([]model.Activation, error) {
	var entities []model.Activation
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r activationRepository) FindById(id string) (model.Activation, error) {
	var entity model.Activation
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r activationRepository) Save(entity model.Activation) (model.Activation, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r activationRepository) Update(entity model.Activation) (model.Activation, error) {
	err := r.DB.Model(model.Activation{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r activationRepository) Delete(entity model.Activation) error {
	return r.DB.Delete(&entity).Error
}

func (r activationRepository) SaveCallback(callback model.ActivationCallback) (model.ActivationCallback, error) {
	err := r.DB.Create(&callback).Error
	return callback, err

}
