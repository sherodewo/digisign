package main

import (
	"io"
	"os"
	"time"

	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/infrastructure/database"
	"los-int-digisign/infrastructure/routes"
	"los-int-digisign/infrastructure/validator"
	"los-int-digisign/model"

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

	// Sentry
	digisign.InitSentry()

	//set credential digisign
	err := digisign.DecryptDigisignCredentials()
	if err != nil {
		tags := map[string]string{
			"app.pkg":     "main",
			"app.func":    "main",
			"app.process": "decrypt-credentials",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		log.Fatal(err.Error())
	}

	//New instance echo
	e := echo.New()

	//Database
	db, err := database.NewDb()
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "main",
			"app.func": "main",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		log.Fatal(err.Error())
	}

	if os.Getenv("APP_ENV") != "production" {
		//Auto migrate
		database.AutoMigrate(db)
		db.Model(&model.Registration{}).AddForeignKey("registration_result_id", "registration_results(id)", "CASCADE", "NO ACTION")
		db.Model(&model.SendDocument{}).AddForeignKey("send_document_result_id", "send_document_results(id)", "CASCADE", "NO ACTION")
		db.Model(&model.Activation{}).AddForeignKey("activation_result_id", "activation_results(id)", "CASCADE", "NO ACTION")
		db.Model(&model.SignDocument{}).AddForeignKey("sign_document_result_id", "sign_document_results(id)", "CASCADE", "NO ACTION")
	}

	// Setup log folder
	if _, err := os.Stat(os.Getenv("LOG_FILE")); os.IsNotExist(err) {
		err = os.MkdirAll(os.Getenv("LOG_FILE"), 0755)
		if err != nil {
			panic(err)
		}
	}
	currentTime := time.Now()

	// Setup Log
	logPath := os.Getenv("LOG_FILE")
	logFileName := currentTime.Format("2006-01-02") + "-" + "digisign.log"
	logFile, err := os.OpenFile(logPath+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "main",
			"app.func": "main",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "LOG_FILE")
		panic("Error create or open log file" + logFileName)
	}

	//Validation
	e.Validator = validator.NewValidator()
	//Set Middleware
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: io.MultiWriter(logFile, os.Stdout),
	}))
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
