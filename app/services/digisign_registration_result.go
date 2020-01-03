package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/models"
)

// DigisignRegistrationResultService :
type DigisignRegistrationResultService struct {
	db *gorm.DB
}

// NewDigisignRegistrationResultService :
func NewDigisignRegistrationResultService(db *gorm.DB) *DigisignRegistrationResultService {
	return &DigisignRegistrationResultService{
		db,
	}
}

// GetByID :
func (us *DigisignRegistrationResultService) GetByID(id string) (*models.DigisignRegistrationResult, error) {
	var m models.DigisignRegistrationResult
	if err := us.db.Where(&models.User{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return &m, nil
}

// Save Result :
func (us *DigisignRegistrationResultService) SaveResult(result string, name bool, id string, notif string, birthplace bool,
	birthdate bool, address string, info string, los_request_id string) (digisignRegistrationResult models.DigisignRegistrationResult, err error) {
	digisignRegistrationResult.ID = id
	digisignRegistrationResult.LosRequestID = los_request_id
	digisignRegistrationResult.Result = result
	digisignRegistrationResult.Notif = notif
	digisignRegistrationResult.Name = name
	digisignRegistrationResult.BirthPlace = birthplace
	digisignRegistrationResult.BirthDate = birthdate
	digisignRegistrationResult.Address = address
	digisignRegistrationResult.Info = info
	err = us.db.Create(&digisignRegistrationResult).Error
	return digisignRegistrationResult, err
}

func (us *DigisignRegistrationResultService) FindById(id string) (c models.DigisignRegistrationResult, err error) {
	if err := us.db.Where("id = ?", id).First(&c).Error; err != nil {
		log.Debug("ERROR", err)
		return c, err
	}
	return c, err
}

func (us *DigisignRegistrationResultService) FindAll() (digisignRegistrationResult []models.DigisignRegistrationResult, err error) {
	err = us.db.Find(&digisignRegistrationResult).Error
	return digisignRegistrationResult, err
}