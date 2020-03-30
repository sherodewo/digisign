package activation

import (
	"encoding/base64"
	"encoding/json"
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/response"
	"kpdigisign/utils"
	"os"
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
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error "+res.String())
	}
	//mapping response
	mapResponse := NewDigisignActivationResponse()
	resMap, err := mapResponse.Bind(dto.ProspectID, res.Body())
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	_, err = c.service.SaveActivation(dto, result, link, res.String())
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", resMap, nil)
}

func (c *Controller) Callback(ctx echo.Context) error {
	encodedValue := ctx.Request().URL.Query().Get("msg")
	decodeValue, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		return response.BadRequest(ctx, utils.BadRequest, nil, err.Error())
	}
	key := os.Getenv("DIGISIGN_AES_KEY")
	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(key))

	var dataMap map[string]interface{}
	err = json.Unmarshal(byteDecrypt, &dataMap)
	if err != nil {
		return response.BadRequest(ctx, utils.BadRequest, nil, err.Error())
	}
	client := NewLosActivationCallbackRequest()
	resLos, err := client.losActivationRequestCallback(dataMap["email"].(string), dataMap["result"].(string),
		dataMap["notif"].(string))

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resLos.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Service callback api Error")
	}
	_, err = c.service.SaveActivationCallback(dataMap["email"].(string), dataMap["result"].(string),
		dataMap["notif"].(string))

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, utils.OK, echo.Map{"message": "Callback success send"}, nil)
}
