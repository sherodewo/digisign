package main

import (
	"kpdigisign/infrastructure/config/digisign"
	"kpdigisign/infrastructure/database"
	"kpdigisign/infrastructure/routes"
	"kpdigisign/infrastructure/validator"
	"kpdigisign/model"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func init() {
	//Load .env file
	err := godotenv.Load("conf/config.env")
	if err != nil {
		panic(err)
	}
}

func main() {

	//set credential digisign
	err := digisign.DecryptDigisignCredentials()
	if err != nil {
		panic(err)
	}

	//New instance echo
	e := echo.New()

	//Database
	db, err := database.NewDb()
	if err != nil {
		panic(err)
	}
	if os.Getenv("APP_ENV") != "production" {
		//Auto migrate
		database.AutoMigrate(db)
		db.Model(&model.Registration{}).AddForeignKey("registration_result_id", "registration_results(id)", "CASCADE", "NO ACTION")
		db.Model(&model.SendDocument{}).AddForeignKey("send_document_result_id", "send_document_results(id)", "CASCADE", "NO ACTION")
		db.Model(&model.Activation{}).AddForeignKey("activation_result_id", "activation_results(id)", "CASCADE", "NO ACTION")
		db.Model(&model.SignDocument{}).AddForeignKey("sign_document_result_id", "sign_document_results(id)", "CASCADE", "NO ACTION")
	}

	//Validation
	e.Validator = validator.NewValidator()
	//Set Middleware
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	//All api route
	routes.ApiRoute(e, db)

	//Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
