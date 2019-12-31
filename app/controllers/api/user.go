package api

import (
	"github.com/labstack/echo"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
	"kpdigisign/app/response"
	"kpdigisign/app/response/mapper"
)

type UserController struct {
	UserRepository repository.UserRepository
}

func (uc *UserController) Store(c echo.Context) error {
	userRequest := request.UserRequest{}
	if err := c.Bind(&userRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	data, err := uc.UserRepository.Create(userRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	serialized := mapper.NewUserMapper().Map(data)
	return response.SingleData(c, "Success Execute request", serialized)
}

func (uc *UserController) GetByID(c echo.Context) error {
	id := c.Param("id")
	data, err := uc.UserRepository.GetByID(id)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	serialized := mapper.NewUserMapper().Map(*data)
	return response.SingleData(c, "Success Execute request", serialized)
}
