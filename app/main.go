package main

import (
	"fmt"
	"log"
	digisignHttpDelivery "los-int-digisign/domain/digisign/delivery/http"
	digisignRepository "los-int-digisign/domain/digisign/repository"
	digisignMultiUsecase "los-int-digisign/domain/digisign/usecase"
	"los-int-digisign/shared/common"
	jsonResponse "los-int-digisign/shared/common/json"
	"los-int-digisign/shared/config"
	"los-int-digisign/shared/database"
	"los-int-digisign/shared/httpclient"
	"los-int-digisign/shared/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"los-int-digisign/docs"

	"github.com/allegro/bigcache"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @contact.name Kredit Plus
// @contact.url https://kreditplus.com
// @contact.email support@kreditplus.com

// @host localhost:9100
// @BasePath /api/v1
// @query.collection.format multi

func main() {

	e := echo.New()

	e.Validator = common.NewValidator()

	config.SetupTimezone()

	config.LoadEnv()

	env := strings.ToLower(config.Env("APP_ENV"))

	config.NewConfiguration(env)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.RequestID())

	e.Debug = config.IsDevelopment

	// digisign, err := database.OpenDigisign()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	los, err := database.OpenLos()
	if err != nil {
		log.Fatal(err)
	}

	logDB, err := database.OpenLosLog()
	if err != nil {
		log.Fatal(err)
	}
	//Set Middleware
	// e.Use(middleware.BodyDumpWithConfig(bodydump.NewBodyDumpMiddleware(logs).BodyDumpConfig()))
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "KREDITPLUS DIGISIGN INTEGRATOR")
	})

	// Cache
	var cache *bigcache.BigCache
	isCacheActive, _ := strconv.ParseBool(config.Env("CACHE_ACTIVE"))
	if isCacheActive {
		cacheExp, _ := strconv.Atoi(config.Env("CACHE_EXPIRED_DEFAULT"))
		cache, _ = bigcache.NewBigCache(bigcache.DefaultConfig(time.Duration(cacheExp) * time.Second))
	}

	utils.NewCache(cache, los, config.IsDevelopment)

	var digisign *gorm.DB

	digiRepo := digisignRepository.NewRepository(digisign, los, logDB)

	httpClient := httpclient.NewHttpClient()

	digiMulti, digiPackage, digiUseCase := digisignMultiUsecase.NewMultiUsecase(digiRepo, httpClient)

	commonJson := jsonResponse.NewResponse()

	apiGroup := e.Group("/api/v1")

	digisignHttpDelivery.DigisignHandler(apiGroup, digiMulti, digiPackage, digiUseCase, commonJson, digiRepo)

	if config.IsDevelopment {
		// programatically set swagger info
		docs.SwaggerInfo.Title = "LOS-INTEGRATION-API"
		docs.SwaggerInfo.Description = "This is a integration api server."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.Env("APP_HOST"), config.Env("APP_PORT"))
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		e.GET("/swagger/*", echoSwagger.WrapHandler)

	} else {

		e.HideBanner = true

		// Newrelic
		app, err := config.InitNewrelic()
		if err == nil {
			e.Use(nrecho.Middleware(app))
		}

	}

	// Setup Server
	e.Start(fmt.Sprintf(":%s", config.Env("APP_PORT")))

}
