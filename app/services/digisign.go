package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
)

// DigisignRegistrationResultService :
type DigisignService struct {
	db *gorm.DB
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
func (us *DigisignService) SaveResult(id string, result string, notif string, jsonResponse string, refTrx string) (digisignResult models.DigisignResult,
	err error) {
	digisignResult.LosID = id
	digisignResult.Result = result
	digisignResult.Notif = notif
	digisignResult.JsonResponse = jsonResponse
	digisignResult.RefTrx = refTrx

	err = us.db.Create(&digisignResult).Error
	return digisignResult, err
}

func (us *DigisignService) SaveDocumentResult(id string, result string, notif string, jsonResponse string,refTrx string) (
	docResult models.DocumentResult, err error) {
	docResult.DocumentID = id
	docResult.Result = result
	docResult.Notif = notif
	docResult.JsonResponse = jsonResponse
	docResult.RefTrx = refTrx

	err = us.db.Create(&docResult).Error
	return docResult, err
}

func (us *DigisignService) SaveDocumentRequest(userId string, documentId string, payment string, sendTo string,
	reqSign string) (doc models.Document, err error) {
	doc.UserID = userId
	doc.DocumentID = documentId
	doc.Payment = payment
	doc.SendTo = sendTo
	doc.ReqSign = reqSign
	err = us.db.Create(&doc).Error
	return doc, err
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
