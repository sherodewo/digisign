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
func NewDigisignController(db *gorm.DB) *api.DigisignController {
	losService:=service.NewLosRequestService(db)
	digisignService:=service.NewDigisignService(db)
	return &api.DigisignController{
		LosRepository:      losService,
		DigisignRepository: digisignService,
	}
}