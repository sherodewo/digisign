package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/helpers"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
)

// UserService :
type UserService struct {
	db *gorm.DB
}

// NewUserService :
func NewUserService(db *gorm.DB) repository.UserRepository {
	return &UserService{
		db,
	}
}

// GetByID :
func (us *UserService) GetByID(id string) (*models.User, error) {
	var m models.User
	if err := us.db.Where(&models.User{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return &m, nil
}

// GetByEmail :
func (us *UserService) GetByEmail(e string) (*models.User, error) {
	var m models.User
	if err := us.db.Where(&models.User{Email: e}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, err
	}
	return &m, nil
}

// Create :
func (us *UserService) Create(u request.UserRequest) (user models.User, err error) {
	passwordHash, err := helpers.HashPassword(u.Password)
	user.Username = u.Username
	user.Email = u.Email
	user.Password = passwordHash
	err = us.db.Create(&user).Error
	return user, err
}

func (us *UserService) Update(id string, um request.UserUpdateRequest) (*models.User, error) {
	data, err := us.FindById(id)
	if err != nil {
		return &data, err
	}
	data.Username = um.Username
	data.Email = um.Email
	data.Password = um.Password

	if err := us.db.Save(&data).Error; err != nil {
		log.Debug("ERROR", err)
		return &data, err
	}
	return &data, nil
}

func (us *UserService) FindById(id string) (c models.User, err error) {
	if err := us.db.Where("id = ?", id).First(&c).Error; err != nil {
		log.Debug("ERROR", err)
		return c, err
	}
	return c, err
}

func (us *UserService) FindAll() (user []models.User, err error) {
	err = us.db.Find(&user).Error
	return user, err
}

func (us *UserService) Destroy(id string) error {
	var m models.User
	if err := us.db.Where(&models.User{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return err
		}
		return err
	}
	if err := us.db.Delete(&models.User{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
