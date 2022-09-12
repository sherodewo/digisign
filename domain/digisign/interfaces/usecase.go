package interfaces

import (
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
)

type Usecase interface {
	DecodeMedia(url string, customerID string) (base64Image string, err error)
	SignUseCase(req request.SignDocDto) (uploadRes response.MediaServiceResponse, err error)
	GeneratePK(prospectID string) (document []byte, docID string, err error)
	SignDocV2(req request.JsonFileSign, prospectID string) (data response.SignDocResponse, err error)
	DownloadDoc(prospectID string, req request.DownloadRequest) (pdfBase64 string, err error)
	UploadDoc(prospectID string, fileName string) (uploadResp response.MediaServiceResponse, err error)
}

type MultiUsecase interface {
	Register(req request.Register) (data response.RegisterResponse, err error)
	ActivationRedirect(msg string) (data interface{}, err error)
}

type Packages interface {
	GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpByte, selfieByte, signatureByte, npwpByte []byte, err error)
	SendDoc(prospectID string) (data response.SendDocResponse, documentID string, err error)
}
