package repository

import (
	"kpdigisign/app/models"
	"kpdigisign/app/request"
)

type DigisignRepository interface {
	GetByID(id string) (*models.DigisignResult, error)
	FindAll() ([]models.DigisignResult, error)
	SaveResult(id string, result string, notif string, jsonResponse string) (models.DigisignResult, error)
	SaveDocumentRequest(request request.LosSendDocumentRequest) (models.Document, error)
	SaveDocumentResult(id string, result string, notif string, jsonResponse string) (models.DocumentResult, error)
}
