package interfaces

import (
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
)

type Usecase interface {
	DecodeMedia(url string, customerID string) (base64Image string, err error)
	// SignUseCase(req request.SignDocDto) (uploadRes response.MediaServiceResponse, err error)
	GeneratePK(prospectID string) (document []byte, docID string, agreementNo string, err error)
	SignDocument(req request.JsonFileSign, prospectID string) (data response.DataSignDocResponse, err error)
	DownloadDoc(prospectID string, req request.DownloadRequest) (pdfBase64 string, err error)
	UploadDoc(prospectID string, fileName string) (uploadResp response.MediaServiceResponse, err error)
	Activation(req request.ActivationRequest) (data response.DataActivationResponse, err error)
}

type MultiUsecase interface {
	Register(req request.Register) (data response.DataRegisterResponse, err error)
	ActivationRedirect(msg string) (data response.DataSignDocResponse, err error)
	SignCallback(msg string) (upload response.MediaServiceResponse, err error)
}

type Packages interface {
	GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpByte, selfieByte, signatureByte, npwpByte []byte, err error)
	SendDoc(req request.SendDoc) (data response.DataSendDocResponse, err error)
}
