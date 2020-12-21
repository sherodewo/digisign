package user

import (
	"github.com/jinzhu/gorm"
	"los-int-digisign/model"
)

type Repository interface {
	FindAll() ([]model.User, error)
	FindById(string) (model.User, error)
	Save(model.User) (model.User, error)
	Update(model.User) (model.User, error)
	Delete(model.User) error
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) Repository {
	return &userRepository{DB: db}
}

func (r userRepository) FindAll() ([]model.User, error) {
	var entities []model.User
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r userRepository) FindById(id string) (model.User, error) {
	var entity model.User
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r userRepository) Save(entity model.User) (model.User, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r userRepository) Update(entity model.User) (model.User, error) {
	err := r.DB.Model(model.User{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r userRepository) Delete(entity model.User) error {
	return r.DB.Delete(&entity).Error
}
