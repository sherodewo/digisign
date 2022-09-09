package repository

import (
	"los-int-digisign/domain/digisign/interfaces"

	"github.com/jinzhu/gorm"
)

type repoHandler struct {
	digisign *gorm.DB
}

func NewRepository(digisign *gorm.DB) interfaces.Repository {
	return &repoHandler{
		digisign: digisign,
	}
}
