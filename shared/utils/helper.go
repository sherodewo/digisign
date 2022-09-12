package utils

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"
)

func GetExtensionImageFromByte(byte []byte) string {

	ext := http.DetectContentType(byte)
	if ext == "text/plain; charset=utf-8" {
		return ""
	}
	res := strings.Split(ext, "image/")

	return res[1]
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

func GenerateTimeNow() string {
	//fetching current time
	currentTime := time.Now().Format(time.RFC3339)
	//differnce between pastdate and current date
	return currentTime
}