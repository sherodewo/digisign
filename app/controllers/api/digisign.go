package api

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
	"kpdigisign/app/client"
	"kpdigisign/app/helpers"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
	"kpdigisign/app/response"
	"kpdigisign/app/response/mapper"
	"os"
)

type DigisignController struct {
	LosRepository      repository.LosRepository
	DigisignRepository repository.DigisignRepository
}

func (d *DigisignController) Register(c echo.Context) error {

	resultMapper := mapper.NewDigisignResultMapper()
	losRequest := request.LosRequest{}
	if err := c.Bind(&losRequest); err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(losRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, helpers.ValidationError, nil, errorData)
	}
	//Check KTP from request
	fileKtp, err := c.FormFile("foto_ktp")
	if fileKtp == nil {
		return response.ValidationError(c, helpers.ValidationError, nil, err)
	}
	bufKtp, err := helpers.GetFileByte("foto_ktp", c)
	//Check Selfie from request
	fileSelfie, err := c.FormFile("foto_selfie")
	if fileSelfie == nil {
		return response.ValidationError(c, helpers.ValidationError, nil, err)
	}
	bufSelfie, err := helpers.GetFileByte("foto_selfie", c)
	//Get NPWP Byte file
	bufNpwp, err := helpers.GetFileByte("foto_npwp", c)
	//Get TTD Byte file
	bufTtd, err := helpers.GetFileByte("tanda_tangan", c)
	//Save request
	data, err := d.LosRepository.Create(&losRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Hit API Registration
	register := client.NewDigisignRegistrationRequest()
	resp, err := register.DigisignRegistration(losRequest.KonsumenType, bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Get Response
	respDigisignRegister := response.NewDigisignResponse()
	if err := respDigisignRegister.Bind(resp.Body()); err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	//Save result register
	resultData, err := d.DigisignRepository.SaveResult(data.ID, respDigisignRegister.JsonFile.Result,
		respDigisignRegister.JsonFile.Notif, resp.String(), respDigisignRegister.JsonFile.RefTrx)

	return response.SingleData(c, helpers.OK, resultMapper.Map(resultData), nil)
}

func (d *DigisignController) SendDocument(c echo.Context) error {
	resultMapper := mapper.NewDocumentResultMapper()
	sendDocRequest := request.LosSendDocumentRequest{}
	if err := c.Bind(&sendDocRequest); err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(sendDocRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, helpers.ValidationError, nil, errorData)
	}
	//Check File Pdf
	file, err := c.FormFile("file")
	if err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	if file == nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err)
	}
	//Save Document Request
	data, err := d.DigisignRepository.SaveDocumentRequest(os.Getenv("DIGISIGN_USER_ID"), sendDocRequest.DocumentID,
		sendDocRequest.Payment, sendDocRequest.SendTo, sendDocRequest.ReqSign)
	if err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	//Get Byte File
	filePdf, err := helpers.GetFileByte("file", c)
	if err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	//Hit API send document
	send := client.NewDigisignSendDocRequest()
	res, err := send.DigisignSendDoc(filePdf, sendDocRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Get Response
	respDigisign := response.NewDigisignResponse()
	if err := respDigisign.Bind(res.Body()); err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())

	}
	//Save Result
	resultData, err := d.DigisignRepository.SaveDocumentResult(data.ID, respDigisign.JsonFile.Result,
		respDigisign.JsonFile.Notif, res.String(), respDigisign.JsonFile.RefTrx)
	return response.SingleData(c, helpers.OK, resultMapper.Map(resultData), nil)
}

func (d DigisignController) Download(c echo.Context) error {

	downloadFileRequest := request.LosDownloadDocumentRequest{}
	if err := c.Bind(&downloadFileRequest); err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(downloadFileRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, helpers.ValidationError, nil, errorData)
	}
	//Hit API download doc
	requestDoc := client.NewDownloadRequest()
	_, file, err := requestDoc.Download(downloadFileRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}

	return response.SingleData(c, helpers.OK, file, nil)
}

func (d DigisignController) DownloadFile(c echo.Context) error {

	downloadFileRequest := request.LosDownloadDocumentRequest{}
	if err := c.Bind(&downloadFileRequest); err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(downloadFileRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, helpers.ValidationError, nil, errorData)
	}
	//Hit API download doc
	requestDoc := client.NewDownloadRequest()
	res, err := requestDoc.DownloadFile(downloadFileRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	return c.Stream(200, "application/pdf", bytes.NewReader(res.Body()))
}

func (d DigisignController) Activation(c echo.Context) error {
	activationRequest := request.LosActivationRequest{}
	if err := c.Bind(&activationRequest); err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(activationRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, helpers.ValidationError, nil, errorData)
	}
	data, err := d.DigisignRepository.SaveActivationRequest(os.Getenv("DIGISIGN_USER_ID"), activationRequest.EmailUser)
	if err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	requestActivation := client.NewActivationRequest()
	_, result, link, err := requestActivation.ActivationDigisign(activationRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	resultActivation, err := d.DigisignRepository.SaveActivationResult(data.ID, result, link)
	return response.SingleData(c, helpers.OK, resultActivation, nil)
}

func (d DigisignController) ActivationCallback(c echo.Context) error {
	query := c.QueryParam("msg")
	log.Info("Query ",query)

	//TODO : MUST DECRPT FROM QUERY PARAM
	decrypt, err := helpers.DecryptAes(query)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//TODO : SAVE CALLBACK
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(*decrypt), &jsonMap)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	result, err := d.DigisignRepository.SaveActivationCallback(jsonMap["email"].(string),jsonMap["result"].(string),
		jsonMap["notif"].(string))

	//TODO : SEND NOTIF TO LOS

	return response.SingleData(c, helpers.OK, result, nil)
}

func (d DigisignController) SignDocument(c echo.Context) error {
	signDocumentRequest := request.LosSignDocumentRequest{}
	if err := c.Bind(&signDocumentRequest); err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	//Validate Request on Struct Los Request
	if err := c.Validate(signDocumentRequest); err != nil {
		errorData := make(echo.Map)
		for _, v := range err.(validator.ValidationErrors) {
			errorData[v.Field()] = v.Tag()
		}
		return response.ValidationError(c, helpers.ValidationError, nil, errorData)
	}
	data, err := d.DigisignRepository.SaveSignDocRequest(os.Getenv("DIGISIGN_USER_ID"), signDocumentRequest.EmailUser,
		signDocumentRequest.DocumentID)
	if err != nil {
		return response.InternalServerError(c, helpers.InternalServerError, nil, err.Error())
	}
	requestSignDoc := client.NewSignDocRequest()
	_, result, link, err := requestSignDoc.DigisignSignDocumentRequest(signDocumentRequest)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	resultSignDocument, err := d.DigisignRepository.SaveSignDocResult(data.ID, result, link)
	return response.SingleData(c, helpers.OK, resultSignDocument, nil)
}

func (d DigisignController) SignDocumentCallback(c echo.Context) error {
	query := c.QueryParam("msg")

	//TODO : MUST DECRPT FROM QUERY PARAM
	decrypt, err := helpers.DecryptAes(query)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	//TODO : SAVE CALLBACK
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(*decrypt), &jsonMap)
	if err != nil {
		return response.BadRequest(c, helpers.BadRequest, nil, err.Error())
	}
	result, err := d.DigisignRepository.SaveSignDocCallback(jsonMap["email"].(string),jsonMap["result"].(string),
		jsonMap["document_id"].(string),jsonMap["status_document"].(string))

	//TODO : SEND NOTIF TO LOS

	return response.SingleData(c, helpers.OK, result, nil)

}
