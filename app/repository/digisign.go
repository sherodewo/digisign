package repository

import (
	"kpdigisign/app/models"
)

type DigisignRepository interface {
	GetByID(id string) (*models.DigisignResult, error)
	FindAll() ([]models.DigisignResult, error)
	SaveResult(id string, result string, notif string, jsonResponse string,refTrx string) (models.DigisignResult, error)
	SaveDocumentRequest(userId string, documentId string, payment string, sendTo string,
		reqSign string) (models.Document, error)
	SaveDocumentResult(id string, result string, notif string, jsonResponse string,refTrx string) (models.DocumentResult, error)
}
