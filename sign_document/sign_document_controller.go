package sign_document

import (
	"encoding/json"
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/response"
	"kpdigisign/utils"
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
	res, result, link, _, err := client.DigisignSignDocumentRequest(dto)
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
	_, err = c.service.SaveSignDocument(dto, result, link, res.String())
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, "Success execute request", resMap, nil)
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
	client := NewLosSignDocumentCallbackRequest()
	resLos, code, message, err := client.losSignDocumentRequestCallback(jsonMap["email"].(string), jsonMap["result"].(string),
		jsonMap["document_id"].(string), jsonMap["status_document"].(string))

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resLos.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Service callback api Error")
	}
	if code != "200" {
		return response.BadRequest(ctx, "Bad Request", nil, "Error hit service callback api")
	}
	_, err = c.service.SaveSignDocumentCallback(jsonMap["document_id"].(string), jsonMap["email"].(string),
		jsonMap["status_document"].(string), jsonMap["result"].(string), )
	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	return response.SingleData(ctx, utils.OK, echo.Map{"code": code, "message": message}, nil)
}
