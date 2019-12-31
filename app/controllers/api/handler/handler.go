package handler

import (
	"github.com/jinzhu/gorm"
	"kpdigisign/app/controllers/api"
	service "kpdigisign/app/services"
)

func NewUserController(db *gorm.DB) *api.UserController {
	userService:=service.NewUserService(db)
	return &api.UserController{UserRepository: userService}
}
