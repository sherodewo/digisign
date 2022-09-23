package sign_document

import (
	"encoding/base64"
	"encoding/json"
	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/infrastructure/response"
	"los-int-digisign/utils"

	"github.com/labstack/echo"
)

type Controller struct {
	service *service
}

func NewSignDocumentController(service *service) Controller {
	return Controller{service: service}
}

func (c *Controller) GetById(ctx echo.Context) error {
	id := ctx.Param("id")
	result, err := c.service.FindSignDocumentById(id)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", result, nil)
}

func (c *Controller) GetAll(ctx echo.Context) error {
	result, err := c.service.FindAllSignDocuments()
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
	client := NewDigisignSignDocumentRequest()

	dto.UserID = digisign.UserID

	res, result, link, _, jsonRequest, err := client.DigisignSignDocumentRequest(dto)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if res.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error "+res.String())
	}
	//mapping response
	mapResponse := NewDigisignSignDocumentResponse()
	resMap, err := mapResponse.Bind(res.Body())
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	_, err = c.service.SaveSignDocument(dto, result, link, res.String(), jsonRequest)
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.MultiData(ctx, "Success execute request", resMap, res.String(), nil)
}

func (c *Controller) Callback(ctx echo.Context) error {
	encodedValue := ctx.Request().URL.Query().Get("msg")
	decodeValue, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		tags := map[string]string{
			"app.pkg":     "sign_document",
			"app.func":    "Callback",
			"app.action":  "decode",
			"app.process": "sign-doc-callback",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DECODE")
		return response.BadRequest(ctx, utils.BadRequest, nil, err.Error())
	}
	key := digisign.AesKey
	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(key))

	var dataMap map[string]interface{}
	err = json.Unmarshal(byteDecrypt, &dataMap)
	if err != nil {
		return response.BadRequest(ctx, utils.BadRequest, nil, err.Error())
	}
	client := NewLosSignDocumentCallbackRequest()
	resLos, err := client.losSignDocumentRequestCallback(dataMap["email_user"].(string), dataMap["result"].(string),
		dataMap["document_id"].(string), dataMap["status_document"].(string))

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resLos.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Service callback api Error")
	}
	_, err = c.service.SaveSignDocumentCallback(dataMap["document_id"].(string), dataMap["email_user"].(string),
		dataMap["status_document"].(string), dataMap["result"].(string), dataMap["notif"].(string))
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.MultiData(ctx, utils.OK, echo.Map{"message": "Callback success send"}, resLos.String(), nil)
}
