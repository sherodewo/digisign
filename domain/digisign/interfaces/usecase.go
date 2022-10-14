package interfaces

import (
	"los-int-digisign/model/entity"
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
)

type Usecase interface {
	DecodeMedia(url string, customerID string) (base64Image string, err error)
	GeneratePK(prospectID string) (document []byte, docID string, agreementNo string, err error)
	DownloadDoc(prospectID string, req request.DownloadRequest) (pdfBase64 string, err error)
	UploadDoc(prospectID string, fileName string) (uploadResp response.MediaServiceResponse, err error)
	CallbackDigisignRegister(data response.DataRegisterResponse, prospectID string) (err error)
	CallbackDigisignActivation(data response.DataActivationResponse, prospectID string) (err error)
	CallbackDigisignSendDoc(data response.DataSendDocResponse, prospectID string) (err error)
	CallbackDigisignSignDoc(data response.DataSignDocResponse, prospectID string) (err error)
	CallbackDigisignSignDocSuccess(prospectID string, url string) (err error)
	DigisignCheck(email, prospectID string) (data response.DataDigisignCheck, err error)
	DownloadAndUpload(prospectID string, req request.DownloadRequest) (uploadResp response.MediaServiceResponse, doc entity.TteDocPk, err error)
}

type MultiUsecase interface {
	Register(req request.Register) (data response.DataRegisterResponse, err error)
	ActivationRedirect(msg string) (data response.DataSignDocResponse, err error)
	SignCallback(msg string) (upload response.MediaServiceResponse, redirectUrl string, err error)
	Activation(req request.ActivationRequest) (data response.DataActivationResponse, err error)
}

type Packages interface {
	GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpByte, selfieByte, signatureByte, npwpByte []byte, err error)
	SendDoc(req request.SendDoc) (data response.DataSendDocResponse, err error)
	SignDocument(req request.JsonFileSign, prospectID string) (data response.DataSignDocResponse, err error)
}
