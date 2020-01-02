package helpers

import (
	"bytes"
	"github.com/labstack/echo"
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
	res := strings.Split(ext, "image/")

	return res[1]
}
