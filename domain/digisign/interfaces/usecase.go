package interfaces

import (
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
)

type Usecase interface {
	DecodeMedia(url string, customerID string) (base64Image string, err error)
	Activation(req request.ActivationRequest) (res response.ActivationResponse, err error)
}

type MultiUsecase interface {
}

type Packages interface {
	GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpBase64, selfieBase64, signatureBase64, npwpBase64 string, err error)
}
