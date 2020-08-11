package digisign

import (
	"kpdigisign/utils"
	"os"
)

var (
	Token  string
	AesKey string
	UserID string
)

func DecryptDigisignCredentials() error {
	decryptDigisignUserID, err := utils.DecryptCredential(os.Getenv("KEY_DECRYPT_CREDENTIALS"), os.Getenv("DIGISIGN_USER_ID"))
	decryptDigisignKey, err := utils.DecryptCredential(os.Getenv("KEY_DECRYPT_CREDENTIALS"), os.Getenv("DIGISIGN_AES_KEY"))
	decryptDigisignToken, err := utils.DecryptCredential(os.Getenv("KEY_DECRYPT_CREDENTIALS"), os.Getenv("DIGISIGN_TOKEN"))
	//set to global var
	UserID = decryptDigisignUserID
	AesKey = decryptDigisignKey
	Token = decryptDigisignToken

	return err
}
