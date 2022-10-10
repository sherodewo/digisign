package httpclient

import (
	"los-int-digisign/model/request"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (_m *MockHttpClient) DigiAPI(url string, method string, param interface{}, file string, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error) {
	args := _m.Called(url, method, param, file, header, timeOut, customerID)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func (_m *MockHttpClient) RegisterAPI(url string, param map[string]string, header map[string]string, method string, timeOut int, dataFile request.DataFile, ProspectID string) (resp *resty.Response, err error) {
	args := _m.Called(url, param, header, method, timeOut, dataFile, ProspectID)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func (_m *MockHttpClient) MediaClient(url string, method string, param map[string]string, file string, header map[string]string, timeOut int, customerID string) (resp *resty.Response, err error) {
	args := _m.Called(url, method, param, file, header, timeOut, customerID)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func (_m *MockHttpClient) SendDocAPI(url, method string, param map[string]string, header map[string]string, timeOut int, dataFile request.DataFile, prospectID string) (resp *resty.Response, err error) {
	args := _m.Called(url, method, param, header, timeOut, dataFile, prospectID)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func (_m *MockHttpClient) SignDocAPI(url, method string, param map[string]string, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error) {
	args := _m.Called(url, method, param, header, timeOut, prospectID)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func (_m *MockHttpClient) ActivationAPI(url, method string, param map[string]string, header map[string]string, timeOut int, prospectID string) (resp *resty.Response, err error) {
	args := _m.Called(url, method, param, header, timeOut, prospectID)
	return args.Get(0).(*resty.Response), args.Error(1)
}
