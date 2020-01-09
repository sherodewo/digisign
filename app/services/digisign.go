package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
)

// DigisignRegistrationResultService :
type DigisignService struct {
	db *gorm.DB
}

func (us *DigisignService) SaveDocumentResult(id string, result string, notif string, jsonResponse string) (
	docResult models.DocumentResult, err error) {
	docResult.DocumentID = id
	docResult.Result = result
	docResult.Notif = notif
	docResult.JsonResponse = jsonResponse

	err = us.db.Create(&docResult).Error
	return docResult, err
}

func (us *DigisignService) SaveDocumentRequest(request request.LosSendDocumentRequest) (doc models.Document, err error) {
	doc.UserID = request.UserID
	doc.DocumentID = request.DocumentID
	doc.Payment = request.Payment
	doc.SendTo = request.SendTo
	doc.ReqSign = request.ReqSign
	err = us.db.Create(&doc).Error
	return doc, err
}

// NewDigisignService :
func NewDigisignService(db *gorm.DB) repository.DigisignRepository {
	return &DigisignService{
		db,
	}
}

// GetByID :
func (us *DigisignService) GetByID(id string) (*models.DigisignResult, error) {
	var m models.DigisignResult
	if err := us.db.Where(&models.User{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return &m, nil
}

// Save Result :
func (us *DigisignService) SaveResult(id string, result string, notif string, jsonResponse string) (digisignResult models.DigisignResult,
	err error) {
	digisignResult.LosID = id
	digisignResult.Result = result
	digisignResult.Notif = notif
	digisignResult.JsonResponse = jsonResponse
	//digisignRegistrationResult.Name = name
	//digisignRegistrationResult.BirthPlace = birthplace
	//digisignRegistrationResult.BirthDate = birthdate
	//digisignRegistrationResult.Address = address
	//digisignRegistrationResult.Info = info
	err = us.db.Create(&digisignResult).Error
	return digisignResult, err
}

func (us *DigisignService) FindById(id string) (c models.DigisignResult, err error) {
	if err := us.db.Where("id = ?", id).First(&c).Error; err != nil {
		log.Debug("ERROR", err)
		return c, err
	}
	return c, err
}

func (us *DigisignService) FindAll() (digisignRegistrationResult []models.DigisignResult, err error) {
	err = us.db.Find(&digisignRegistrationResult).Error
	return digisignRegistrationResult, err
}
