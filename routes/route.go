package routes

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/controllers/api/handler"
	"kpdigisign/config"
)

func New() (e *echo.Echo) {
	e = echo.New()
	//Route group
	api:=e.Group("/api")
	v1:=api.Group("/v1")
	//DB
	db := config.New()
	//Validation
	e.Validator = NewValidator()

	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))


	userController:=handler.NewUserController(db)
	v1.POST("/user", userController.Store)
	v1.GET("/user/:id", userController.GetByID)

	return e
}
