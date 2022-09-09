package httpclient

import (
	"encoding/json"
	"errors"
	"los-int-digisign/model/response"
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
	MediaClient(url, method string, param interface{}, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error)
}

func (h httpClient) MediaAPI(url string, param interface{}, header map[string]string, method string, timeOut int, retry bool, countRetry interface{}) (resp *resty.Response, err error) {

	header["Authorization"] = os.Getenv("AUTH_MEDIA")

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
	case "POST":
		resp, err = client.R().SetHeaders(header).SetBody(param).Post(url)
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return
}

func (h httpClient) MediaClient(url, method string, param interface{}, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error) {

	client := resty.New()

	client.SetTimeout(time.Second * time.Duration(timeOut))

	switch method {
	case "GET":
		resp, err = client.R().SetHeaders(header).Get(url)
	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	var mediaResponse response.MediaServiceResponse

	json.Unmarshal(resp.Body(), &mediaResponse)

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

	case "POST":
		resp, err = client.R().SetHeaders(header).SetBody(param).Post(url)
	case "GET":
		resp, err = client.R().SetHeaders(header).SetBody(param).Get(url)
	case "PUT":
		resp, err = client.R().SetHeaders(header).SetBody(param).Put(url)
	case "DELETE":
		resp, err = client.R().SetHeaders(header).SetBody(param).Delete(url)

	}

	if err != nil {
		err = errors.New("connection error")
		return
	}

	return

}
