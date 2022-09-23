package response

import (
	"github.com/labstack/echo"
	"net/http"
)

func SingleData(c echo.Context, message string, data interface{}, error interface{}) error {
	return c.JSON(http.StatusOK, Single{
		Meta: Meta{
			Code:    http.StatusOK,
			Message: message,
			Error:   error,
		},
		Data: data,
	})
}

func MultiData(c echo.Context, message string, data interface{}, respOri interface{}, error interface{}) error {
	return c.JSON(http.StatusOK, Single{
		Meta: Meta{
			Code:    http.StatusOK,
			Message: message,
			Error:   error,
			RespOri: respOri,
		},
		Data: data,
	})
}

func NotFound(c echo.Context, message string, data interface{}, error interface{}) error {
	return c.JSON(http.StatusNotFound, Single{
		Meta: Meta{
			Code:    http.StatusNotFound,
			Message: message,
			Error:   error,
		},
		Data: data,
	})
}

func BadRequest(c echo.Context, message string, data interface{}, error interface{}) error {
	return c.JSON(http.StatusBadRequest, Single{
		Meta: Meta{
			Code:    http.StatusBadRequest,
			Message: message,
			Error:   error,
		},
		Data: data,
	})
}

func ValidationError(c echo.Context, message string, data interface{}, error interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, Single{
		Meta: Meta{
			Code:    http.StatusUnprocessableEntity,
			Message: message,
			Error:   error,
		},
		Data: data,
	})
}

func InternalServerError(c echo.Context, message string, data interface{}, error interface{}) error {
	return c.JSON(http.StatusInternalServerError, Single{
		Meta: Meta{
			Code:    http.StatusInternalServerError,
			Message: message,
			Error:   error,
		},
		Data: data,
	})
}

func Unauthorized(c echo.Context, message string, data interface{}, error interface{}) error {
	return c.JSON(http.StatusUnauthorized, Single{
		Meta: Meta{
			Code:    http.StatusUnauthorized,
			Message: message,
			Error:   error,
		},
		Data: data,
	})
}

