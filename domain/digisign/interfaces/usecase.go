package interfaces

import (
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
)

type Usecase interface {
	DecodeMedia(url string, customerID string) (base64Image string, err error)
	SignUseCase(req request.SignDocDto) (uploadRes response.MediaServiceResponse, err error)
	Activation(req request.ActivationRequest) (res response.ActivationResponse, err error)
}

type MultiUsecase interface {
	Register(req request.Register) (data response.RegisterResponse, err error)
	ActivationRedirect(msg string) (data interface{}, err error)
}

type Packages interface {
	GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpByte, selfieByte, signatureByte, npwpByte []byte, err error)
	SendDoc(req request.SendDocRequest) (data response.SendDocResponse, err error)
}
