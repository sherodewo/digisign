package api

import (
	"github.com/labstack/echo"
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
	var bufNpwp, bufTtd []byte

	//Check KTP
	fileKtp, err := c.FormFile("foto_ktp")
	if fileKtp == nil {
		return response.BadRequest(c, "NOT FOUND KTP", nil)
	}
	bufKtp, err := helpers.GetFileByte("foto_ktp", c)
	//Check Selfie
	fileSelfie, err := c.FormFile("foto_selfie")
	if fileSelfie == nil {
		return response.BadRequest(c, "NOT FOUND Selfie", nil)
	}
	bufSelfie, err := helpers.GetFileByte("foto_selfie", c)
	//Get NPWP
	bufNpwp, err = helpers.GetFileByte("foto_npwp", c)
	//Get TTD
	bufTtd, err = helpers.GetFileByte("tanda_tangan", c)

	data, err := d.LosRepository.Create(losRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	register := client.NewDigisignRegistrationRequest()
	resp, err := register.DigisignRegistration(losRequest.KonsumenType, bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}

	respDigisignRegister := response.NewDigisignResponse()
	if err := respDigisignRegister.Bind(resp.Body()); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	resultData, err := d.DigisignRepository.SaveResult(data.ID, respDigisignRegister.JsonFile.Result,
		respDigisignRegister.JsonFile.Notif, resp.String())

	return response.SingleData(c, "Success execute resuest", resultMapper.Map(resultData))
}

func (d *DigisignController) SendDocument(c echo.Context) error {
	resultMapper := mapper.NewDocumentResultMapper()

	sendDocRequest := request.LosSendDocumentRequest{}
	if err := c.Bind(&sendDocRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Check File Pdf
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}
	if file == nil {
		return response.BadRequest(c, "NOT FOUND File", nil)
	}
	//===============
	//Save Document Request
	data, err := d.DigisignRepository.SaveDocumentRequest(sendDocRequest.UserId,sendDocRequest.DocumentId,
		sendDocRequest.Payment,sendDocRequest.SendTo,sendDocRequest.ReqSign)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}
	//=====================

	//Get
	filePdf, err := helpers.GetFileByte("file", c)
	send := client.NewDigisignSendDocRequest()
	res, err := send.DigisignSendDoc(filePdf, c.FormValue("userId"),c.FormValue("documentId"),sendDocRequest)

	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}

	respDigisignRegister := response.NewDigisignResponse()
	if err := respDigisignRegister.Bind(res.Body()); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	resultData, err := d.DigisignRepository.SaveDocumentResult(data.ID, respDigisignRegister.JsonFile.Result,
		respDigisignRegister.JsonFile.Notif, res.String())
	return response.SingleData(c, "Success execute resuest", resultMapper.Map(resultData))

}

func (d DigisignController) Download(c echo.Context) error  {

	downloadFileRequest := request.LosDownloadDocumentRequest{}
	if err := c.Bind(&downloadFileRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}

	requestDoc:=client.NewDownloadRequest()
	res,err:=requestDoc.Download(downloadFileRequest)
	if err != nil{
		return response.BadRequest(c, err.Error(), nil)
	}

	_, err = c.Response().Write(res.Body())

	return nil
	//return response.SingleData(c, "Success execute resuest", resultMapper.Map(resultData))
}