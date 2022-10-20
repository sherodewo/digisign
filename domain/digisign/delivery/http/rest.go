package http

import (
	"fmt"
	"los-int-digisign/domain/digisign/interfaces"
	"los-int-digisign/model/request"
	"los-int-digisign/shared/common"
	"los-int-digisign/shared/constant"
	"os"

	"github.com/labstack/echo/v4"
)

type digisignHandler struct {
	multiUsecase interfaces.MultiUsecase
	packages     interfaces.Packages
	usecase      interfaces.Usecase
	repository   interfaces.Repository
	Json         common.JSON
}

func DigisignHandler(route *echo.Group, multiUsecase interfaces.MultiUsecase, packages interfaces.Packages, usecase interfaces.Usecase, json common.JSON, repository interfaces.Repository) {
	handler := digisignHandler{
		multiUsecase: multiUsecase,
		packages:     packages,
		usecase:      usecase,
		repository:   repository,
		Json:         json,
	}

	digiGroup := route.Group("/digisign")
	{
		digiGroup.POST("/register", handler.Register)
		digiGroup.POST("/sign-doc", handler.SignDoc)
		digiGroup.GET("/activation/callback", handler.ActivationCallback)
		digiGroup.GET("/sign-document/callback", handler.SignCallback)
		digiGroup.POST("/activation", handler.Activation)
		digiGroup.POST("/send-doc", handler.SendDoc)
		digiGroup.POST("/step-check", handler.CheckDigisignIndex)
		digiGroup.POST("/download/base64", handler.DownloadDoc)
	}
}

// Digisign godoc
// @Description Api Register Digisign
// @Tags Digisign
// @Produce json
// @Param body body request.Register true "Body payload"
// @Success 200 {object} response.Api{}
// @Failure 400 {object} response.Api{error=response.ErrorValidation}
// @Failure 500 {object} response.Api{}
// @Router /digisign/register [post]
func (h *digisignHandler) Register(ctx echo.Context) (err error) {

	var req request.Register

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS Digisign - Register ", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS Digisign - Register", err)
	}

	data, err := h.multiUsecase.Register(req)

	if err != nil {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_timeout - Register Timeout"))
	}

	return h.Json.Ok(ctx, "LOS Digisign - Register", data)

}

// Digisign godoc
// @Description Api Activation Digisign
// @Tags Digisign
// @Produce json
// @Param body body request.ActivationRequest true "Body payload"
// @Success 200 {object} response.Api{}
// @Failure 400 {object} response.Api{error=response.ErrorValidation}
// @Failure 500 {object} response.Api{}
// @Router /digisign/activation [post]
func (h *digisignHandler) Activation(ctx echo.Context) (err error) {

	var req request.ActivationRequest

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS Digisign - Activation", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS Digisign - Activation", err)
	}

	data, err := h.multiUsecase.Activation(req)
	if err != nil {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_timeout - Activation Timeout"))
	}

	return h.Json.Ok(ctx, "LOS Digisign", data)
}

// Digisign godoc
// @Description Api Send Doc Digisign
// @Tags Digisign
// @Produce json
// @Param body body request.SendDoc true "Body payload"
// @Success 200 {object} response.Api{}
// @Failure 400 {object} response.Api{error=response.ErrorValidation}
// @Failure 500 {object} response.Api{}
// @Router /digisign/send-doc [post]
func (h *digisignHandler) SendDoc(ctx echo.Context) (err error) {

	var req request.SendDoc

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS Digisign - SendDoc", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS Digisign - SendDoc", err)
	}

	data, err := h.packages.SendDoc(req)
	if err != nil {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_timeout - Send Doc Timeout"))
	}

	return h.Json.Ok(ctx, "LOS Digisign", data)
}

func (h *digisignHandler) ActivationCallback(ctx echo.Context) (err error) {

	msg := ctx.QueryParam("msg")

	data, err := h.multiUsecase.ActivationRedirect(msg)

	if err != nil && err.Error() != constant.EXPIRED {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_error - Activation Redirect Error"))
	}

	data.Link = "www.google.com"
	if data.Link == "" {
		return h.Json.Ok(ctx, "LOS Digisign", data)
	}

	if err != nil {
		return h.Json.Ok(ctx, "LOS Digisign - Activation Callback", "OK")
	}

	return h.Json.Ok(ctx, "CALLBACK ACTIVATION" ,data)
}

// Digisign godoc
// @Description Api sign Doc Digisign
// @Tags Digisign
// @Produce json
// @Param body body request.SignDocRequest true "Body payload"
// @Success 200 {object} response.Api{}
// @Failure 400 {object} response.Api{error=response.ErrorValidation}
// @Failure 500 {object} response.Api{}
// @Router /digisign/sign-doc [post]
func (h *digisignHandler) SignDoc(ctx echo.Context) (err error) {

	var req request.SignDocRequest

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS - Sign Document", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS - Sign Document", err)
	}

	sign, err := h.packages.SignDocument(request.JsonFileSign{
		UserID:     os.Getenv("DIGISIGN_USER_ID"),
		DocumentID: req.DocumentID,
		Email:      req.Email,
	}, req.ProspectID)

	if err != nil {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_error - Sign Document Error"))
	}

	return h.Json.Ok(ctx, "LOS - Sign Document", sign)
}

// SignCallback godoc
// @Description Api SignCallback
// @Tags Digisign
// @Produce json
// @Param msg query string true
// @Success 200 {object} response.Api{}
// @Failure 400 {object} response.Api{error=response.ErrorValidation}
// @Failure 500 {object} response.Api{}
// @Router /digisign/sign-document/callback [get]
func (h *digisignHandler) SignCallback(ctx echo.Context) (err error) {

	msg := ctx.QueryParam("msg")

	_, redirect, err := h.multiUsecase.SignCallback(msg)

	if err != nil && err.Error() != constant.EXPIRED {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_error - Activation Redirect Error"))
	}

	if err != nil {
		return h.Json.Ok(ctx, "LOS Digisign - Sign Callback", "OK")
	}

	return h.Json.Ok(ctx, "CALLBACK ACTIVATION" , echo.Map{"redirect_url":redirect})
}

// Digisign godoc
// @Description Api Check Step Digisign
// @Tags Digisign
// @Produce json
// @Param body body request.DigisignCheck true "Body payload"
// @Success 200 {object} response.Api{}
// @Failure 400 {object} response.Api{error=response.ErrorValidation}
// @Failure 500 {object} response.Api{}
// @Router /digisign/step-check [post]
func (h *digisignHandler) CheckDigisignIndex(ctx echo.Context) (err error) {

	var req request.DigisignCheck

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS - Digisign Check", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS - Digisign Check", err)
	}

	data, err := h.usecase.DigisignCheck(req.Email, req.ProspectID)

	if err != nil {
		return h.Json.NotFound(ctx, "LOS - Digisign Check")

	}

	return h.Json.Ok(ctx, "LOS - Digisign Check", data)

}

func (h *digisignHandler) DownloadDoc(ctx echo.Context) (err error) {

	var req request.DigisignDownload

	if err := ctx.Bind(&req); err != nil {
		return h.Json.InternalServerError(ctx, "LOS - Download Document", err)
	}

	if err := ctx.Validate(&req); err != nil {
		return h.Json.BadRequestErrorValidation(ctx, "LOS - Download Document", err)
	}

	_, err = h.usecase.DownloadDoc(req.ProspectID, request.DownloadRequest{
		UserID:     os.Getenv("DIGISIGN_USER_ID"),
		DocumentID: req.DocumentID,
	})

	if err != nil {
		return h.Json.ServerSideError(ctx, "LOS Digisign", fmt.Errorf("upstream_service_error - Download Document Error"))
	}

	return h.Json.Ok(ctx, "LOS - Download Document", "download document success")
}

func (h *digisignHandler) TestRedirect(ctx echo.Context) (err error) {

	return ctx.Redirect(200, "http://www.google.com")
}
