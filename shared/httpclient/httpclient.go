package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
	"los-int-digisign/shared/config"
	"los-int-digisign/shared/constant"
	"los-int-digisign/shared/utils"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type httpClient struct{}

func NewHttpClient() HttpClient {
	return &httpClient{}
}

type HttpClient interface {
	MediaAPI(url string, param interface{}, header map[string]string, method string, timeOut int, retry bool, countRetry interface{}) (resp *resty.Response, err error)
	EngineAPI(url string, param interface{}, header map[string]string, method string, timeOut int, retry bool, countRetry interface{}) (resp *resty.Response, err error)
	ActivationAPI(url, method string, param map[string]string, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error)
	RegisterAPI(url string, param map[string]string, header map[string]string, method string, timeOut int, dataFile request.DataFile, ProspectID string) (resp *resty.Response, err error)
	MediaClient(url, method string, param map[string]string, file string, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error)
	SendDocAPI(url, method string, param map[string]string, header map[string]string, timeOut int, dataFile request.DataFile, prospectID string) (resp *resty.Response, err error)
	SignDocAPI(url, method string, param map[string]string, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error)
	DigiAPI(url, method string, param map[string]string, file string, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error)
	DocumentAPI(url, method string, param interface{}, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error)
}

func (h httpClient) MediaAPI(url string, param interface{}, header map[string]string, method string, timeOut int, retry bool, countRetry interface{}) (resp *resty.Response, err error) {

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}
	if retry {
		client.SetRetryCount(countRetry.(int))
		client.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= 500
			})
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {
	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetBody(param).Post(url)
	}

	if err != nil {
		err = errors.New(constant.CONNECTION_ERROR)
		return
	}

	return
}

func (h httpClient) MediaClient(url, method string, param map[string]string, file string, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error) {
	client := resty.New()

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {
	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetFormData(param).SetFile("file", file).Post(url)
	case constant.METHOD_GET:
		resp, err = client.R().SetHeaders(header).Get(url)
	}

	if err != nil {
		err = errors.New(constant.CONNECTION_ERROR)
		return
	}

	var mediaResponse response.MediaServiceResponse

	json.Unmarshal(resp.Body(), &mediaResponse)

	return

}

func (h httpClient) RegisterAPI(url string, param map[string]string, header map[string]string, method string, timeOut int, dataFile request.DataFile, ProspectID string) (resp *resty.Response, err error) {

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {

	case constant.METHOD_POST:
		if dataFile.Signature == nil && dataFile.PhotoNPWP == nil {
			resp, err = client.R().SetHeaders(header).
				SetFormData(param).SetFileReader("fotoktp", "ktp_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.PhotoKTP), bytes.NewReader(dataFile.PhotoKTP)).
				SetFileReader("fotodiri", "selfie_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.Selfie), bytes.NewReader(dataFile.Selfie)).
				Post(url)
		} else if dataFile.PhotoNPWP == nil {
			resp, err = client.R().SetHeaders(header).
				SetFormData(param).SetFileReader("fotoktp", "ktp_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.PhotoKTP), bytes.NewReader(dataFile.PhotoKTP)).
				SetFileReader("fotodiri", "selfie_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.Selfie), bytes.NewReader(dataFile.Selfie)).
				SetFileReader("ttd", "ttd_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.Signature), bytes.NewReader(dataFile.Signature)).
				Post(url)
		} else if dataFile.Signature == nil {
			resp, err = client.R().SetHeaders(header).
				SetFormData(param).SetFileReader("fotoktp", "ktp_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.PhotoKTP), bytes.NewReader(dataFile.PhotoKTP)).
				SetFileReader("fotodiri", "selfie_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.Selfie), bytes.NewReader(dataFile.Selfie)).
				SetFileReader("fotonpwp", "npwp_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.PhotoNPWP), bytes.NewReader(dataFile.PhotoNPWP)).
				Post(url)
		} else {
			resp, err = client.R().SetHeaders(header).
				SetFormData(param).SetFileReader("fotoktp", "ktp_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.PhotoKTP), bytes.NewReader(dataFile.PhotoKTP)).
				SetFileReader("fotodiri", "selfie_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.Selfie), bytes.NewReader(dataFile.Selfie)).
				SetFileReader("ttd", "ttd_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.Signature), bytes.NewReader(dataFile.Signature)).
				SetFileReader("fotonpwp", "npwp_"+dataFile.Name+"."+utils.GetExtensionImageFromByte(dataFile.PhotoNPWP), bytes.NewReader(dataFile.PhotoNPWP)).
				Post(url)
		}

	}

	var registerResponse response.RegisterResponse
	json.Unmarshal(resp.Body(), &registerResponse)

	isError := false

	if resp.StatusCode() != 200 {
		isError = true
	}

	logs := map[string]interface{}{
		"ID":            ProspectID,
		"response":      registerResponse,
		"response_code": resp.StatusCode(),
		"url":           url,
		"response_time": fmt.Sprintf("%dms", resp.Time().Milliseconds()),
	}

	go config.SetCustomLog("API_DIGISIGN", isError, logs, "REGISTER-API")

	if err != nil {
		err = errors.New(constant.CONNECTION_ERROR)
		return
	}

	return

}

func (h httpClient) EngineAPI(url string, param interface{}, header map[string]string, method string, timeOut int, retry bool, countRetry interface{}) (resp *resty.Response, err error) {

	header["Content-Type"] = "application/json"

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}
	if retry {
		client.SetRetryCount(countRetry.(int))
		client.AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return r.StatusCode() >= 500
			})
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {

	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetBody(param).Post(url)
	case constant.METHOD_GET:
		resp, err = client.R().SetHeaders(header).SetBody(param).Get(url)

	}

	if err != nil {
		err = errors.New(constant.CONNECTION_ERROR)
		return
	}

	return

}

func (h httpClient) SendDocAPI(url, method string, param map[string]string, header map[string]string, timeOut int, dataFile request.DataFile, prospectID string) (resp *resty.Response, err error) {

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {

	case constant.METHOD_POST:
		if dataFile.DocumentPK != nil {
			resp, err = client.R().SetHeaders(header).
				SetFormData(param).SetFileReader("file", "file_"+dataFile.DocumentID+".pdf", bytes.NewReader(dataFile.DocumentPK)).
				Post(url)
		} else {
			return nil, fmt.Errorf("connection error")
		}
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return

}

func (h httpClient) SignDocAPI(url, method string, param map[string]string, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error) {

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {

	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetFormData(param).Post(url)
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return

}

func (h httpClient) ActivationAPI(url, method string, param map[string]string, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error) {

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {

	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetFormData(param).Post(url)
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return

}

func (h httpClient) DigiAPI(url, method string, param map[string]string, file string, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error) {

	client := resty.New()
	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {
	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetFormData(param).Post(url)
	case constant.METHOD_GET:
		resp, err = client.R().SetHeaders(header).Get(url)
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return

}

func (h httpClient) DocumentAPI(url, method string, param interface{}, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error) {

	client := resty.New()

	if os.Getenv("APP_ENV") != "production" {
		client.SetDebug(true)
	}

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {

	case constant.METHOD_POST:
		resp, err = client.R().SetHeaders(header).SetBody(param).Post(url)
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return
}
