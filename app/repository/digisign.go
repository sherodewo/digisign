package repository

import (
	"kpdigisign/app/models"
)

type DigisignRepository interface {
	GetByID(id string) (*models.DigisignResult, error)
	FindAll() ([]models.DigisignResult, error)
	SaveResult(id string, result string, notif string, jsonResponse string, refTrx string) (models.DigisignResult, error)
	SaveDocumentRequest(userId string, documentId string, payment string, sendTo string,
		reqSign string) (models.Document, error)
	SaveDocumentResult(id string, result string, notif string, jsonResponse string, refTrx string) (
		models.DocumentResult, error)
	SaveActivationRequest(userId string, emailUser string) (models.Activation, error)
	SaveActivationResult(activationId string, result string, link string) (models.ActivationResult, error)
	SaveActivationCallback(email string, result string, notif string) (models.ActivationCallback, error)
	SaveSignDocRequest(userId string, emailUser string, documentId string) (models.SignDocument, error)
	SaveSignDocResult(signDocumentId string, result string, link string) (models.SignDocumentResult, error)
	SaveSignDocCallback(email string, result string, documentId string, statusDocument string) (
		models.SignDocumentCallback, error)
}
