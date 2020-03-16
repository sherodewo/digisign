package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetFileByte(key string, c echo.Context) (byte []byte, err error) {
	// Source
	file, err := c.FormFile(key)
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, src)

	return buf.Bytes(), nil
}

func GetExtensionImageFromByte(byte []byte) string {

	ext := http.DetectContentType(byte)
	log.Info("DATA EXT = ", ext)
	if ext == "text/plain; charset=utf-8" {
		return ""
	}
	res := strings.Split(ext, "image/")

	return res[1]
}

func CheckDocumentFile(key string, c echo.Context) (code int, err error) {
	file, err := c.FormFile(key)
	if err != nil {
		return 0, err
	}
	if file == nil {
		return 404, nil
	}
	return 200, nil
}

func DecryptAes(message string) (*string, error) {
	key := []byte(os.Getenv("DIGISIGN_AES_KEY"))
	cipherText, err := base64.URLEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short!")
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	decodedmess := string(cipherText)
	return &decodedmess, nil
}
