package json

import (
	"encoding/json"
	"fmt"
	"los-int-digisign/shared/common"
	"los-int-digisign/shared/utils"
	"net/http"
	"strings"

	resp "los-int-digisign/model/response"

	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type (
	response struct {
	}
)

func NewResponse() common.JSON {
	return &response{}
}

func (c *response) Ok(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusOK, resp.Api{
		Message:    message,
		Errors:     nil,
		Data:       data,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) InternalErrorWithMessage(ctx echo.Context, message string, err error) error {
	return ctx.JSON(http.StatusInternalServerError, resp.Api{
		Message:    message,
		Errors:     err.Error(),
		Data:       nil,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) OkCreated(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusCreated, resp.Api{
		Message:    message,
		Errors:     nil,
		Data:       data,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) ServiceUnavailable(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusServiceUnavailable, resp.Api{
		Message:    message,
		Errors:     "service_unavailable",
		Data:       nil,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) InternalServerError(ctx echo.Context, message string, err error) error {
	errString := handleInternalError(err)
	return ctx.JSON(http.StatusInternalServerError, resp.Api{
		Message:    message + " - " + errString,
		Errors:     "internal_server_error",
		Data:       nil,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) NotFound(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusNotFound, resp.Api{
		Message:    message,
		Data:       nil,
		Errors:     "not_found",
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) BadRequestErrorValidation(ctx echo.Context, message string, err error) error {
	var errors = make([]resp.ErrorValidation, len(err.(validator.ValidationErrors)))

	for k, v := range err.(validator.ValidationErrors) {
		errors[k] = resp.ErrorValidation{
			Field:   strcase.ToSnake(v.Field()),
			Message: formatMessage(v),
		}
	}
	return ctx.JSON(http.StatusBadRequest, resp.Api{
		Message:    message,
		Errors:     errors,
		Data:       nil,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func (c *response) ServerSideError(ctx echo.Context, message string, err error) error {

	handleError := strings.Split(err.Error(), " - ")

	switch handleError[0] {

	case "upstream_service_error":
		return ctx.JSON(http.StatusBadGateway, resp.Api{
			Message:    fmt.Sprintf("%s - %s", message, handleError[1]),
			Errors:     "upstream_service_error",
			Data:       nil,
			ServerTime: utils.GenerateTimeNow(),
		})

	case "upstream_service_timeout":
		return ctx.JSON(http.StatusGatewayTimeout, resp.Api{
			Message:    fmt.Sprintf("%s - %s", message, handleError[1]),
			Errors:     "upstream_service_timeout",
			Data:       nil,
			ServerTime: utils.GenerateTimeNow(),
		})

	case "service_unavailable":
		return ctx.JSON(http.StatusServiceUnavailable, resp.Api{
			Message:    fmt.Sprintf("%s - %s", message, handleError[1]),
			Errors:     "service_unavailable",
			Data:       nil,
			ServerTime: utils.GenerateTimeNow(),
		})

	case "bad_request":
		return ctx.JSON(http.StatusBadRequest, resp.Api{
			Message:    fmt.Sprintf("%s - %s", message, handleError[1]),
			Errors:     "bad_request",
			Data:       nil,
			ServerTime: utils.GenerateTimeNow(),
		})
	}

	return err
}

func (c *response) Unauthorized(ctx echo.Context, message string, err error) error {
	return ctx.JSON(http.StatusUnauthorized, resp.Api{
		Message:    message,
		Errors:     err.Error(),
		Data:       nil,
		ServerTime: utils.GenerateTimeNow(),
	})
}

func handleUnmarshalError(err error) []resp.ErrorValidation {
	var apiErrors []resp.ErrorValidation
	if he, ok := err.(*echo.HTTPError); ok {
		if ute, ok := he.Internal.(*json.UnmarshalTypeError); ok {
			valError := resp.ErrorValidation{
				Field:   ute.Field,
				Message: ute.Error(),
			}
			apiErrors = append(apiErrors, valError)
		}
		if se, ok := he.Internal.(*json.SyntaxError); ok {
			valError := resp.ErrorValidation{
				Field:   "Syntax Error",
				Message: se.Error(),
			}
			apiErrors = append(apiErrors, valError)
		}
		if iue, ok := he.Internal.(*json.InvalidUnmarshalError); ok {
			valError := resp.ErrorValidation{
				Field:   iue.Type.String(),
				Message: iue.Error(),
			}
			apiErrors = append(apiErrors, valError)
		}
	}
	return apiErrors
}

func handleInternalError(err error) (apiErrors string) {

	if he, ok := err.(*echo.HTTPError); ok {
		if _, ok := he.Internal.(*json.UnmarshalTypeError); ok {
			apiErrors = "Unmarshal Type Error"
			return
		}
		if _, ok := he.Internal.(*json.SyntaxError); ok {
			apiErrors = "Syntax Error"
			return
		}
		if _, ok := he.Internal.(*json.InvalidUnmarshalError); ok {
			apiErrors = "Invalid Unmarshal Error"
			return
		}

		if strings.Contains(err.Error(), "unexpected EOF") {
			apiErrors = "Unexpected EOF"
			return
		}

		if strings.Contains(err.Error(), "unexpected end") {
			apiErrors = "Unexpected end Of JSON Input"
			return
		}

	}

	apiErrors = "Other"
	return
}

func formatMessage(err validator.FieldError) string {

	_ = err.Param()

	message := fmt.Sprintf("Field validation for '%s' failed on the '%s'", strcase.ToSnake(err.Field()), err.Tag())

	//switch err.Tag() {
	//case constant.TAG_GT:
	//	message = fmt.Sprintf("accepted:gt=%s", param)
	//}
	return message
}
