package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
)

// DigisignRegistrationResultService :
type DigisignResultService struct {
	db *gorm.DB
}

// NewDigisignRegistrationResultService :
func NewDigisignRegistrationResultService(db *gorm.DB) repository.DigisignResultRepository {
	return &DigisignResultService{
		db,
	}
}

// GetByID :
func (us *DigisignResultService) GetByID(id string) (*models.DigisignResult, error) {
	var m models.DigisignResult
	if err := us.db.Where(&models.User{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return &m, nil
}

// Save Result :
func (us *DigisignResultService) SaveResult(result string, notif string, los_id string) (digisignResult models.DigisignResult, err error) {
	digisignResult.LosID = los_id
	digisignResult.Result = result
	digisignResult.Notif = notif
	//digisignRegistrationResult.Name = name
	//digisignRegistrationResult.BirthPlace = birthplace
	//digisignRegistrationResult.BirthDate = birthdate
	//digisignRegistrationResult.Address = address
	//digisignRegistrationResult.Info = info
	err = us.db.Create(&digisignResult).Error
	return digisignResult, err
}

func (us *DigisignResultService) FindById(id string) (c models.DigisignResult, err error) {
	if err := us.db.Where("id = ?", id).First(&c).Error; err != nil {
		log.Debug("ERROR", err)
		return c, err
	}
	return c, err
}

func (us *DigisignResultService) FindAll() (digisignRegistrationResult []models.DigisignResult, err error) {
	err = us.db.Find(&digisignRegistrationResult).Error
	return digisignRegistrationResult, err
}