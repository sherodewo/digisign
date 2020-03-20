package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
)

func GetExtensionImageFromByte(byte []byte) string {

	ext := http.DetectContentType(byte)
	if ext == "text/plain; charset=utf-8" {
		return ""
	}
	res := strings.Split(ext, "image/")

	return res[1]
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

	decodedMessage := string(cipherText)
	return &decodedMessage, nil
}

func Base64Decode(data string) ([]byte, error) {
	if data == "" {
		return nil, errors.New("String null")
	}
	decode, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return decode, nil
}


