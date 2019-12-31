package request

import (
	"github.com/labstack/echo"
)

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cr *UserRequest) Bind(c echo.Context) (*UserRequest, error) {
	if err := c.Bind(cr); err != nil {
		return nil, err
	}
	return cr, nil
}
