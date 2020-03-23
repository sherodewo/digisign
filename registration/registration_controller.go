package registration

import (
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/response"
	"kpdigisign/utils"
)

type Controller struct {
	service *service
}

func NewRegistrationController(service *service) Controller {
	return Controller{service: service}
}

func (c *Controller) GetById(ctx echo.Context) error {
	id := ctx.Param("id")
	result, err := c.service.FindRegistrationById(id)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", result, nil)
}

func (c *Controller) GetAll(ctx echo.Context) error {
	result, err := c.service.FindAllRegistrations()
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", result, nil)
}

func (c *Controller) Store(ctx echo.Context) error {
	var dto Dto
	if err := ctx.Bind(&dto); err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if err := ctx.Validate(&dto); err != nil {
		return response.ValidationError(ctx, "Validation error", nil, err.Error())
	}

	//required file
	bufktp, err := utils.Base64Decode(dto.FotoKtp)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	bufSelfie, err := utils.Base64Decode(dto.FotoSelfie)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	//optional file
	bufNpwp, _ := utils.Base64Decode(dto.FotoNpwp)
	bufTtd, _ := utils.Base64Decode(dto.FotoTandaTangan)

	client := NewDigisignRegistrationRequest()
	resp, result, notif, reftrx, jsonResponse, err := client.DigisignRegistration(dto.KonsumenType, bufktp, bufSelfie,
		bufNpwp, bufTtd, dto)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resp.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error "+resp.String())
	}
	mapResponse := NewDigisignRegistrationResponse()
	resMap, err := mapResponse.Bind(resp.Body())
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	_, err = c.service.SaveRegistration(dto, result, notif, reftrx, jsonResponse)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", resMap, nil)
}
