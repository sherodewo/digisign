package repository

import (
	"kpdigisign/app/models"
)

type DigisignResultRepository interface {
	GetByID(id string) (*models.DigisignResult, error)
	FindAll() ([]models.DigisignResult, error)
	SaveResult(result string, notif string, los_id string)(models.DigisignResult, error)
}
