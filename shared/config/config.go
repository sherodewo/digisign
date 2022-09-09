package config

import (
	"io"
	"los-int-digisign/shared/constant"
	"los-int-digisign/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	logger "github.com/labstack/gommon/log"
	"github.com/newrelic/go-agent/v3/newrelic"
	log "github.com/sirupsen/logrus"
)

var (
	DateLogFile   map[string]string
	GetLogFile    map[string]*os.File
	IsDevelopment bool
)

func LoadEnv() {
	err := godotenv.Load("conf/config.env")
	if err != nil {
		log.Fatal("Error loading env file")
	}
}

func NewConfiguration(appEnv string) {

	if strings.ToLower(appEnv) != "prod" && strings.ToLower(appEnv) != "production" {
		IsDevelopment = true
	} else {
		IsDevelopment = false
	}

	GetLogFile = make(map[string]*os.File)
	DateLogFile = make(map[string]string)

}

func Env(key string) string {
	env, err := godotenv.Read("conf/config.env")
	if err != nil {
		logger.Fatalf("Error %v", err)
	}
	v := env[key]
	return v
}

func InitNewrelic() (*newrelic.Application, error) {
	newrelicActive, _ := strconv.ParseBool(Env("NEWRELIC_ACTIVE"))
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(Env("APP_NAME")),
		newrelic.ConfigLicense(Env("NEWRELIC_CONFIG_LICENSE")),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(addConfig *newrelic.Config) {
			addConfig.Enabled = newrelicActive
			addConfig.Labels = map[string]string{
				"Env": strings.ToLower(Env("APP_ENV")),
				"Tag": strings.ToLower(Env("APP_VERSION")),
			}
		},
	)
	return app, err
}

func CreateCustomLogFile(keyConfig string) {

	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(loc)

	// create folder general
	logPath := Env("LOG_FILE") + keyConfig + "/"

	active, _ := strconv.ParseBool(Env(keyConfig))
	if logPath != "" && active {
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			err = os.MkdirAll(logPath, 0755)
			if err != nil {
				panic(err)
			}
		}

		logFileName := strings.ToLower(keyConfig) + "-" + currentTime.Format(constant.TIME_FORMAT) + ".log"
		logFile, err := os.OpenFile(logPath+logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			logger.Fatalf("error opening file: %v", err)
		}
		GetLogFile[keyConfig] = logFile
		DateLogFile[keyConfig] = currentTime.Format(constant.TIME_FORMAT)
	}
}

func SetCustomLog(keyConfig string, isError bool, data map[string]interface{}, msg string) {

	dateNow := time.Now().Format("2006-01-02")
	if DateLogFile[keyConfig] != dateNow {
		CreateCustomLogFile(keyConfig)
	}
	logPath := os.Getenv("LOG_FILE")
	active, _ := strconv.ParseBool(os.Getenv(keyConfig))
	if logPath != "" && active {

		logFile := GetLogFile[keyConfig]

		log.SetOutput(io.MultiWriter(logFile, os.Stdout))
		log.SetFormatter(&log.JSONFormatter{})
		if isError {
			log.WithFields(data).Error(msg)
			return
		}
		log.WithFields(data).Info(msg)
		return
	}
}

func DigisignDBCredential() (string, string, string, string, string) {
	user, _ := utils.DecryptCredential(os.Getenv("DIGISIGN_DB_USERNAME"))
	pwd, _ := utils.DecryptCredential(os.Getenv("DIGISIGN_DB_PASSWORD"))
	host, _ := utils.DecryptCredential(os.Getenv("DIGISIGN_DB_HOST"))
	port, _ := utils.DecryptCredential(os.Getenv("DIGISIGN_DB_PORT"))
	database, _ := utils.DecryptCredential(os.Getenv("DIGISIGN_DB_DATABASE"))

	return user, pwd, host, port, database
}
