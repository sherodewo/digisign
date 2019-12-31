package repository

import (
	"kpdigisign/app/models"
	"kpdigisign/app/request"
)

type UserRepository interface {
	GetByID(id string) (*models.User, error)
	FindAll() ([]models.User, error)
	GetByEmail(string) (*models.User, error)
	Update(string, request.UserUpdateRequest) (*models.User, error)
	Destroy(string) error
	Create(request request.UserRequest) (models.User, error)
}
