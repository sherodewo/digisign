package activation

import (
	"encoding/json"
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/response"
	"kpdigisign/utils"
)

type Controller struct {
	service *service
}

func NewActivationController(service *service) Controller {
	return Controller{service: service}
}

func (c *Controller) GetById(ctx echo.Context) error {
	id := ctx.Param("id")
	result, err := c.service.FindActivationById(id)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", result, nil)
}

func (c *Controller) GetAll(ctx echo.Context) error {
	result, err := c.service.FindAllActivations()
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
	client := NewDigisignActivationRequest()
	res, result, link, err := client.ActivationDigisign(dto)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if res.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error")
	}
	data, err := c.service.SaveActivation(dto, result, link)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", data, nil)
}

func (c *Controller) Callback(ctx echo.Context) error {
	query := ctx.QueryParam("msg")

	decrypt, err := utils.DecryptAes(query)
	if err != nil {
		return response.BadRequest(ctx, utils.BadRequest, nil, err.Error())
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(*decrypt), &jsonMap)

	if err != nil {
		return response.BadRequest(ctx, utils.BadRequest, nil, err.Error())
	}
	client := NewLosActivationCallbackRequest()
	resLos, code, message, err := client.losActivationRequestCallback(jsonMap["email"].(string), jsonMap["result"].(string),
		jsonMap["notif"].(string))

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resLos.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Service callback api Error")
	}
	if code !="200" {
		return response.BadRequest(ctx, "Bad Request", nil, "Error hit service callback api")
	}
	_, err = c.service.SaveActivationCallback(jsonMap["email"].(string), jsonMap["result"].(string),
		jsonMap["notif"].(string))

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, utils.OK, echo.Map{"code":code,"message":message}, nil)
}
