package repository

import (
	"kpdigisign/app/models"
	"kpdigisign/app/request"
)

type LosRequestRepository interface {
	GetByID(id string) (*models.LosRequest, error)
	FindAll() ([]models.LosRequest, error)
	GetByEmail(string) (*models.LosRequest, error)
	Destroy(string) error
	Create(request request.LosRequest) (models.LosRequest, error)
	SaveResult(result string,info string, emailRegistered string, name bool, birthplace bool,
		birthdate bool, Addres string, selfieMatch bool ) (models.DigisignRegistrationResult, error)
}
