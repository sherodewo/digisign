package routes

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/controllers/api/handler"
	"kpdigisign/app/models"
	"kpdigisign/config"
)

func New() (e *echo.Echo) {
	e = echo.New()
	e.Static("/static", "assets")
	//Route group
	api := e.Group("/api")
	v1 := api.Group("/v1")
	//DB
	db := config.New()
	config.AutoMigrate(db)
	db.Model(&models.DigisignResult{}).AddForeignKey("los_id", "los(id)", "CASCADE", "NO ACTION")
	db.Model(&models.DocumentResult{}).AddForeignKey("document_id", "documents(id)", "CASCADE", "NO ACTION")

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

	userController := handler.NewUserController(db)
	v1.POST("/user", userController.Store)
	v1.GET("/user/:id", userController.GetByID)

	digisignController := handler.NewDigisignController(db)
	v1.POST("/digisign/register", digisignController.Register)
	v1.POST("/digisign/send-document", digisignController.SendDocument)
	v1.POST("/digisign/download", digisignController.Download)
	v1.POST("/digisign/download/file", digisignController.DownloadFile)
	v1.POST("/digisign/activation", digisignController.Activation)
	v1.GET("/digisign/activation/callback", digisignController.ActivationCallback)
	v1.POST("/digisign/sign-document", digisignController.SignDocument)
	v1.GET("/digisign/sign-document/callback", digisignController.SignDocumentCallback)

	return e
}
