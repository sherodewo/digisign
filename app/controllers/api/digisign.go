package api

import (
	"bytes"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
	"kpdigisign/app/client"
	"kpdigisign/app/helpers"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
	"kpdigisign/app/response"
	"kpdigisign/app/response/mapper"
)

type DigisignController struct {
	LosRepository      repository.LosRepository
	DigisignRepository repository.DigisignRepository
}

func (d *DigisignController) Register(c echo.Context) error {

	resultMapper := mapper.NewDigisignResultMapper()
	losRequest := request.LosRequest{}
	if err := c.Bind(&losRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(losRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, "Validation Error", errorData)
	}
	//Check KTP from request
	fileKtp, err := c.FormFile("foto_ktp")
	if fileKtp == nil {
		return response.ValidationError(c, "NOT FOUND KTP", nil)
	}
	bufKtp, err := helpers.GetFileByte("foto_ktp", c)
	//Check Selfie from request
	fileSelfie, err := c.FormFile("foto_selfie")
	if fileSelfie == nil {
		return response.ValidationError(c, "NOT FOUND Selfie", nil)
	}
	bufSelfie, err := helpers.GetFileByte("foto_selfie", c)
	//Get NPWP Byte file
	bufNpwp, err := helpers.GetFileByte("foto_npwp", c)
	//Get TTD Byte file
	bufTtd, err := helpers.GetFileByte("tanda_tangan", c)
	//Save request
	data, err := d.LosRepository.Create(&losRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Hit API Registration
	register := client.NewDigisignRegistrationRequest()
	resp, err := register.DigisignRegistration(losRequest.KonsumenType, bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}
	//Get Response
	respDigisignRegister := response.NewDigisignResponse()
	if err := respDigisignRegister.Bind(resp.Body()); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}
	//Save result register
	resultData, err := d.DigisignRepository.SaveResult(data.ID, respDigisignRegister.JsonFile.Result,
		respDigisignRegister.JsonFile.Notif, resp.String(), respDigisignRegister.JsonFile.RefTrx)

	return response.SingleData(c, "Success execute resuest", resultMapper.Map(resultData))
}

func (d *DigisignController) SendDocument(c echo.Context) error {
	resultMapper := mapper.NewDocumentResultMapper()
	sendDocRequest := request.LosSendDocumentRequest{}
	if err := c.Bind(&sendDocRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(sendDocRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, "Validation Error", errorData)
	}
	//Check File Pdf
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}
	if file == nil {
		return response.BadRequest(c, "NOT FOUND File", nil)
	}
	//Save Document Request
	data, err := d.DigisignRepository.SaveDocumentRequest(sendDocRequest.UserID, sendDocRequest.DocumentID,
		sendDocRequest.Payment, sendDocRequest.SendTo, sendDocRequest.ReqSign)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}
	//Get Byte File
	filePdf, err := helpers.GetFileByte("file", c)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}
	//Hit API send document
	send := client.NewDigisignSendDocRequest()
	res, err := send.DigisignSendDoc(filePdf, sendDocRequest)
	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}
	//Get Response
	respDigisign := response.NewDigisignResponse()
	if err := respDigisign.Bind(res.Body()); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}
	//Save Result
	resultData, err := d.DigisignRepository.SaveDocumentResult(data.ID, respDigisign.JsonFile.Result,
		respDigisign.JsonFile.Notif, res.String(), respDigisign.JsonFile.RefTrx)
	return response.SingleData(c, "Success execute resuest", resultMapper.Map(resultData))
}

func (d DigisignController) Download(c echo.Context) error {

	downloadFileRequest := request.LosDownloadDocumentRequest{}
	if err := c.Bind(&downloadFileRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(downloadFileRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, "Validation Error", errorData)
	}
	//Hit API download doc
	requestDoc := client.NewDownloadRequest()
	_, file, err := requestDoc.Download(downloadFileRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	return response.SingleData(c, "Success execute resuest", file)
}

func (d DigisignController) DownloadFile(c echo.Context) error {

	downloadFileRequest := request.LosDownloadDocumentRequest{}
	if err := c.Bind(&downloadFileRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(downloadFileRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, "Validation Error", errorData)
	}
	//Hit API download doc
	requestDoc := client.NewDownloadRequest()
	res, err := requestDoc.DownloadFile(downloadFileRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return c.Stream(200, "application/pdf", bytes.NewReader(res.Body()))
}
