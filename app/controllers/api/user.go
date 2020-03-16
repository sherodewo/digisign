package api

import (
	"github.com/labstack/echo"
	"kpdigisign/app/helpers"
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
		return response.BadRequest(c, helpers.BadRequest, nil,err.Error())
	}
	data, err := uc.UserRepository.Create(userRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil,err.Error())

	}
	serialized := mapper.NewUserMapper().Map(data)
	return response.SingleData(c, helpers.OK, serialized,nil)
}

func (uc *UserController) GetByID(c echo.Context) error {
	id := c.Param("id")
	data, err := uc.UserRepository.GetByID(id)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil,err.Error())
	}
	serialized := mapper.NewUserMapper().Map(*data)
	return response.SingleData(c, helpers.OK, serialized,nil)

}
