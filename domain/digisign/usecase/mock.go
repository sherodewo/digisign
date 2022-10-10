package usecase

import (
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"

	"github.com/stretchr/testify/mock"
)

type MockUsecase struct {
	mock.Mock
}

type MockPackages struct {
	mock.Mock
}

type MockMultiUsecase struct {
	mock.Mock
}

func (_m *MockUsecase) DecodeMedia(url string, customerID string) (base64Image string, err error) {
	args := _m.Called(url, customerID)
	return args.String(0), args.Error(1)
}

func (_m *MockUsecase) SignUseCase(req interface{}) (uploadRes interface{}, err error) {
	args := _m.Called(req)
	return args.Get(0), args.Error(1)
}

func (_m *MockMultiUsecase) Register(req interface{}) (data interface{}, err error) {
	args := _m.Called(req)
	return args.Get(0), args.Error(1)
}

func (_m *MockMultiUsecase) ActivationRedirect(msg string) (data interface{}, err error) {
	args := _m.Called(msg)
	return args.Get(0), args.Error(1)
}

func (_m *MockPackages) GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpByte, selfieByte, signatureByte, npwpByte []byte, err error) {
	args := _m.Called(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID)
	return args.Get(0).([]byte), args.Get(1).([]byte), args.Get(2).([]byte), args.Get(3).([]byte), args.Error(4)
}

func (_m *MockPackages) SendDoc(req request.SendDocRequest) (data response.SendDocResponse, err error) {
	args := _m.Called(req)
	return args.Get(0).(response.SendDocResponse), args.Error(1)
}
