package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"kpdigisign/infrastructure/config"
	"net/http"
)

func ApiRoute(e *echo.Echo, db *gorm.DB) {

	//Route group
	api := e.Group("/api")
	v1 := api.Group("/v1")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "KREDITPLUS DIGISIGN API")
	})

	//Route user
	userController := config.InjectUserController(db)
	v1.GET("/users", userController.GetAll)
	v1.GET("/user/:id", userController.GetById)
	v1.POST("/user", userController.Store)
	v1.PUT("/user/:id", userController.Update)
	v1.DELETE("/user/:id", userController.Delete)

	//Route digisign
	digisign := v1.Group("/digisign")

	//registration
	registrationController := config.InjectRegistrationController(db)
	digisign.GET("/registrations", registrationController.GetAll)
	digisign.GET("/registration/:id", registrationController.GetById)
	digisign.POST("/registration", registrationController.Store)

	//send document
	sendDocumentController := config.InjectSendDocumentController(db)
	digisign.GET("/send-documents", sendDocumentController.GetAll)
	digisign.GET("/send-document/:id", sendDocumentController.GetById)
	digisign.POST("/send-document", sendDocumentController.Store)

	//download document
	downloadDocumentController := config.InjectDownloadDocumentController()
	digisign.GET("/document/download/file", downloadDocumentController.DownloadDocument)
	digisign.GET("/document/download/base64", downloadDocumentController.DownloadDocumentBase64)

	//activation
	activationController := config.InjectActivationController(db)
	digisign.GET("/activations", activationController.GetAll)
	digisign.GET("/activation/:id", activationController.GetById)
	digisign.POST("/activation", activationController.Store)
	digisign.GET("/activation/callback", activationController.Callback)

	//sign document
	signDocumentController := config.InjectSignDocumentController(db)
	digisign.GET("/sign-documents", signDocumentController.GetAll)
	digisign.GET("/sign-document/:id", signDocumentController.GetById)
	digisign.POST("/sign-document", signDocumentController.Store)
	digisign.GET("/sign-document/callback", signDocumentController.Callback)

}
