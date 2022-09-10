package interfaces

import "los-int-digisign/model/request"

type Usecase interface {
	DecodeMedia(url string, customerID string) (base64Image string, err error)
	SignDoc(req request.SignDocRequest) (err error)
}

type MultiUsecase interface {
}

type Packages interface {
	GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpBase64, selfieBase64, signatureBase64, npwpBase64 string, err error)
}
