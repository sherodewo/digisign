package send_document

import (
	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/infrastructure/response"
	"los-int-digisign/utils"

	"github.com/labstack/echo"
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
		tags := map[string]string{
			"app.pkg":    "send_document",
			"app.func":   "Store",
			"app.action": "decode",
			"app.process": "base64File",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DECODE")
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	client := NewDigisignSendDocRequest()

	dto.UserID = digisign.UserID

	res, result, notif, reftrx, jsonResponse, jsonRequest, err := client.DigisignSendDoc(byteFile, dto)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if res.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error "+res.String())
	}
	data, err := c.service.SaveSendDocument(dto, result, notif, reftrx, jsonResponse, jsonRequest)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", data, nil)
}
