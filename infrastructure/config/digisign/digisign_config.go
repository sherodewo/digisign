package digisign

import (
	"errors"
	"fmt"
	"los-int-digisign/utils"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/gommon/log"
)

var (
	Token  string
	AesKey string
	UserID string
)

func DecryptDigisignCredentials() error {
	decryptDigisignUserID, err := utils.DecryptCredential(os.Getenv("DIGISIGN_USER_ID"))
	decryptDigisignKey, err := utils.DecryptCredential(os.Getenv("DIGISIGN_AES_KEY"))
	decryptDigisignToken, err := utils.DecryptCredential(os.Getenv("DIGISIGN_TOKEN"))
	//set to global var
	UserID = decryptDigisignUserID
	AesKey = decryptDigisignKey
	Token = decryptDigisignToken

	return err
}

func InitSentry() {
	isSentryActive, errEnv := strconv.ParseBool(os.Getenv("SENTRY_ACTIVE"))
	if isSentryActive && errEnv == nil {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:         os.Getenv("SENTRY_DSN"),
			Environment: os.Getenv("APP_ENV"),
		}); err != nil {
			tags := map[string]string{
				"app.pkg":  "config",
				"app.func": "InitSentry",
			}

			SendToSentry(tags, nil, "SENTRY")
			log.Info("Error : ", err)
		}
		defer sentry.Flush(2 * time.Second)
	} else {
		errEnv = errors.New("Init Sentry")
		log.Info("Error : ", errEnv)
	}
}

func SendToSentry(tags map[string]string, extras map[string]interface{}, errCategory interface{}) {
	isSentryActive, errEnv := strconv.ParseBool(os.Getenv("SENTRY_ACTIVE"))
	if isSentryActive && errEnv == nil {
		sentry.ConfigureScope(func(scope *sentry.Scope) {

			//scope.SetTag(tagKey, tagVal)
			if len(tags) > 0 {
				for k, v := range tags {
					scope.SetTag(k, v)
				}
			}

			if len(extras) > 0 {
				for k, v := range extras {
					scope.SetExtra(k, v)
				}
			}

			title := fmt.Sprintf("%s | %s", os.Getenv("APP_NAME"), errCategory)
			sentry.CaptureMessage(title)

		})
	}
}
