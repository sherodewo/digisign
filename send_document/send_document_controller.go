package send_document

import (
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/response"
	"kpdigisign/utils"
)

type Controller struct {
	service *service
}

func NewSendDocumentController(service *service) Controller {
	return Controller{service: service}
}

func (c *Controller) GetById(ctx echo.Context) error {
	id := ctx.Param("id")
	result, err := c.service.FindSendDocumentById(id)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", result, nil)
}

func (c *Controller) GetAll(ctx echo.Context) error {
	result, err := c.service.FindAllSendDocuments()
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
	byteFile, err := utils.Base64Decode(dto.File)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	client := NewDigisignSendDocRequest()
	res,result, notif, reftrx, jsonResponse, err := client.DigisignSendDoc(byteFile, dto)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if res.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil,"Digisign api error")
	}
	data, err := c.service.SaveSendDocument(dto, result, notif, reftrx, jsonResponse)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", data, nil)
}
