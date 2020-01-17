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

func (us *DigisignService) SaveDocumentResult(id string, result string, notif string, jsonResponse string, refTrx string) (
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

func (us *DigisignService) SaveActivationRequest(userId string, emailUser string) (activation models.Activation, err error) {
	activation.UserID = userId
	activation.EmailUser = emailUser
	err = us.db.Create(&activation).Error
	return activation, err
}

func (us *DigisignService) SaveActivationResult(activationId string, result string, link string) (
	activationResult models.ActivationResult, err error) {
	activationResult.ActivationID = activationId
	activationResult.Result = result
	activationResult.Link = link
	err = us.db.Create(&activationResult).Error
	return activationResult, err
}

func (us *DigisignService) SaveActivationCallback(email string, result string, notif string) (
	activationCallback models.ActivationCallback, err error) {
	activationCallback.Result = result
	activationCallback.Notif = notif
	activationCallback.Email = email
	err = us.db.Create(&activationCallback).Error
	return activationCallback, err
}

func (us *DigisignService) SaveSignDocRequest(userId string, emailUser string, documentId string) (
	signDocRequest models.SignDocument, err error) {
	signDocRequest.EmailUser = emailUser
	signDocRequest.UserID = userId
	signDocRequest.DocumentID = documentId
	err = us.db.Create(&signDocRequest).Error
	return signDocRequest, err
}

func (us *DigisignService) SaveSignDocResult(signDocumentId string, result string, link string) (
	signDocResult models.SignDocumentResult, err error) {
	signDocResult.SignDocumentID = signDocumentId
	signDocResult.Result = result
	signDocResult.Link = link
	err = us.db.Create(&signDocResult).Error
	return signDocResult, err
}

func (us *DigisignService) SaveSignDocCallback(email string, result string, documentId string, statusDocument string) (
	signDocCallback models.SignDocumentCallback, err error) {
	signDocCallback.DocumentID = documentId
	signDocCallback.Email = email
	signDocCallback.Result = result
	signDocCallback.StatusDocument = statusDocument

	err = us.db.Create(&signDocCallback).Error
	return signDocCallback, err
}
