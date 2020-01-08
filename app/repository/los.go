package repository

import (
	"kpdigisign/app/models"
	"kpdigisign/app/request"
)

type LosRepository interface {
	GetByID(id string) (*models.Los, error)
	FindAll() ([]models.Los, error)
	GetByEmail(string) (*models.Los, error)
	Destroy(string) error
	Create(request request.LosRequest) (models.Los, error)
	SaveResult(id string,result string,notif string,jsonResponse string ) (models.DigisignResult, error)
}
