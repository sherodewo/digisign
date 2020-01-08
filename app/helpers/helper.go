package helpers

import (
	"bytes"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"strings"
)

func GetImageByte(key string, c echo.Context) (byte []byte, err error) {
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
