package repository

import (
	"kpdigisign/app/models"
)

type DigisignRegistrationResultRepository interface {
	GetByID(id string) (*models.DigisignRegistrationResult, error)
	FindAll() ([]models.DigisignRegistrationResult, error)
	SaveResult(result string, name bool, id string, notif string, birthplace bool,
				birthdate bool, address string, info string, los_request_id string)(models.DigisignRegistrationResult, error)
}
