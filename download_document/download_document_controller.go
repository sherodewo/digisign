package download_document

import (
	"bytes"
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/response"
)

type Controller struct{}

func NewDownloadDocumentController() Controller {
	return Controller{}
}

func (c *Controller) DownloadDocument(ctx echo.Context) error {
	var dto Dto
	if err := ctx.Bind(&dto); err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if err := ctx.Validate(&dto); err != nil {
		return response.ValidationError(ctx, "Validation Error", nil, err.Error())
	}
	client := NewDigisignDownloadRequest()
	resp, err := client.DownloadFile(dto)

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resp.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error")
	}
	return ctx.Stream(200, "application/pdf", bytes.NewReader(resp.Body()))

}
func (c *Controller) DownloadDocumentBase64(ctx echo.Context) error {
	var dto Dto
	if err := ctx.Bind(&dto); err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if err := ctx.Validate(&dto); err != nil {
		return response.ValidationError(ctx, "Validation Error", nil, err.Error())
	}
	client := NewDigisignDownloadRequest()
	resp, file, err := client.DownloadFileBase64(dto)

	if err != nil {
		return response.BadRequest(ctx, "Bad Request", nil, err.Error())
	}
	if resp.IsError() {
		return response.BadRequest(ctx, "Bad Request", nil, "Digisign api error")
	}
	return response.SingleData(ctx, "Success execute request", file, nil)

}
