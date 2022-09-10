package http

import (
	"github.com/labstack/echo/v4"
	"los-int-digisign/domain/digisign/interfaces"
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
	"los-int-digisign/shared/common"
)

type digisignHandler struct {
	multiUsecase interfaces.MultiUsecase
	usecase      interfaces.Usecase
	repository   interfaces.Repository
	Json         common.JSON
}

func DigisignHandler(route *echo.Group, multiUsecase interfaces.MultiUsecase, usecase interfaces.Usecase, json common.JSON, repository interfaces.Repository) {
	handler := digisignHandler{
		multiUsecase: multiUsecase,
		usecase:      usecase,
		repository:   repository,
		Json:         json,
	}

	digiGroup := route.Group("/document")
	{
		digiGroup.POST("/sign", handler.SignDoc)
	}
}

func (h *digisignHandler) SignDoc(ctx echo.Context) (err error) {
	var req request.SignDocDto

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS - Bind Sign Doc", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS - Validate Sign Doc", err)
	}

	sign, err := h.usecase.SignUseCase(req)

	resp := response.SignResponse{
		ProspectID: req.ProspectID,
		Url:        sign.Data.MediaURL,
	}
	return h.Json.Ok(ctx, "SUCCESS", resp)
}
