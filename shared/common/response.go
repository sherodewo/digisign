package common

import (
	"github.com/labstack/echo/v4"
)

type JSON interface {
	Ok(ctx echo.Context, message string, data interface{}) error
	OkCreated(ctx echo.Context, message string, data interface{}) error
	ServiceUnavailable(ctx echo.Context, message string) error
	InternalServerError(ctx echo.Context, message string, err error) error
	BadRequestErrorValidation(ctx echo.Context, message string, err error) error
	ServerSideError(ctx echo.Context, message string, err error) error
	Unauthorized(ctx echo.Context, message string, err error) error
	NotFound(ctx echo.Context, message string) error
	InternalErrorWithMessage(ctx echo.Context, message string, err error) error
}
