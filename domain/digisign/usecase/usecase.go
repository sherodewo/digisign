package usecase

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"los-int-digisign/domain/digisign/interfaces"

	"los-int-digisign/model/entity"
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
	"los-int-digisign/shared/config"
	"los-int-digisign/shared/constant"
	"los-int-digisign/shared/httpclient"
	"los-int-digisign/shared/utils"

	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type (
	packages struct {
		repository interfaces.Repository
		usecase    interfaces.Usecase
		httpclient httpclient.HttpClient
	}
	usecase struct {
		repository interfaces.Repository
		httpclient httpclient.HttpClient
	}
	multiUsecase struct {
		repository interfaces.Repository
		httpclient httpclient.HttpClient
		usecase    interfaces.Usecase
		packages   interfaces.Packages
	}
)

func NewPackages(repository interfaces.Repository, usecase interfaces.Usecase, httpclient httpclient.HttpClient) interfaces.Packages {
	return &packages{repository: repository, usecase: usecase, httpclient: httpclient}
}

func NewUsecase(repository interfaces.Repository, httpclient httpclient.HttpClient) interfaces.Usecase {
	return &usecase{
		repository: repository,
		httpclient: httpclient,
	}
}

func NewMultiUsecase(repository interfaces.Repository, httpclient httpclient.HttpClient) (interfaces.MultiUsecase, interfaces.Packages, interfaces.Usecase) {

	usecase := NewUsecase(repository, httpclient)
	packages := NewPackages(repository, usecase, httpclient)

	return &multiUsecase{
		repository: repository,
		httpclient: httpclient,
		usecase:    usecase,
		packages:   packages,
	}, packages, usecase
}

func (u multiUsecase) Register(req request.Register) (data response.DataRegisterResponse, err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_95S"))

	registerParam := request.RegisterRequest{
		UserID:     os.Getenv("DIGISIGN_USER_ID"),
		Address:    req.Address,
		Gender:     req.Gender,
		Kecamatan:  req.Kecamatan,
		Kelurahan:  req.Kelurahan,
		City:       req.City,
		Name:       req.Name,
		Phone:      req.Phone,
		TglLahir:   req.BirthDate,
		Provinci:   "INDONESIA",
		IDKtp:      req.IDKtp,
		BirthPlace: req.BirthPlace,
		Email:      req.Email,
		NPWP:       req.NPWP,
		Redirect:   true,
	}

	ktpByte, selfieByte, signatureByte, npwpByte, err := u.packages.GetRegisterPhoto(req.PhotoKTP, req.Selfie, req.Signature, req.PhotoNPWP, req.ProspectID)

	if err != nil {
		return
	}

	dataFile := request.DataFile{
		PhotoKTP:  ktpByte,
		Selfie:    selfieByte,
		Signature: signatureByte,
		PhotoNPWP: npwpByte,
		Name:      req.Name,
	}

	jsonField, _ := json.Marshal(request.JsonFile{
		JsonFile: registerParam,
	})

	param := map[string]string{
		"jsonfield": string(jsonField),
	}

	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": "Bearer " + os.Getenv("DIGISIGN_TOKEN"),
	}

	var digiResp response.RegisterResponse

	dummyActive, _ := strconv.ParseBool(os.Getenv("DUMMY_ACTIVE"))

	if dummyActive {

		var dummy entity.DigisignDummy

		dummy, err = u.repository.GetDigisignDummy(req.Email, "REGISTER")

		json.Unmarshal([]byte(dummy.Response), &digiResp)

		go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
			ProspectID: req.ProspectID,
			Response:   dummy.Response,
			Activity:   "REGISTER",
		})

		data = RegisterMappingResponse(digiResp, req.ProspectID)

		return
	}

	resp, err := u.httpclient.RegisterAPI(os.Getenv("REGISTER_URL"), param, header, constant.METHOD_POST, timeOut, dataFile, req.ProspectID)

	if err != nil {
		return
	}

	go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
		ProspectID: req.ProspectID,
		Response:   string(resp.Body()),
		Activity:   "REGISTER",
	})

	json.Unmarshal(resp.Body(), &digiResp)

	isError := false

	if resp.StatusCode() != 200 {
		isError = true
	}

	logs := map[string]interface{}{
		"ID":            req.ProspectID,
		"request":       param,
		"response":      digiResp,
		"response_code": resp.StatusCode(),
		"url":           os.Getenv("REGISTER_URL"),
		"response_time": fmt.Sprintf("%dms", resp.Time().Milliseconds()),
	}

	go config.SetCustomLog("API_DIGISIGN", isError, logs, "REGISTER-API")

	data = RegisterMappingResponse(digiResp, req.ProspectID)

	go u.usecase.CallbackDigisignRegister(data, req.ProspectID)

	return
}

func (u usecase) CallbackDigisignRegister(data response.DataRegisterResponse, prospectID string) (err error) {

	if data.Decision == constant.DECISION_REJECT {

		dataWorker, _ := u.repository.GetDataWorker(prospectID)

		if dataWorker.TransactionType != "" {
			timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

			header := map[string]string{
				"Content-Type":  "application/json",
				"Authorization": os.Getenv("AUTH_KPM"),
			}

			headerParam, _ := json.Marshal(header)

			param, _ := json.Marshal(map[string]interface{}{
				"prospect_id": prospectID,
				"code":        data.Code,
				"status":      constant.REGISTER_FAILED,
				"activity":    constant.ACTIVITY_FINISHED,
			})

			paramCallbackStatus, _ := json.Marshal(map[string]interface{}{
				"messages": "LOS DigiSign",
				"errors":   nil,
				"data": map[string]interface{}{
					"prospect_id":     prospectID,
					"code":            data.Code,
					"decision":        constant.DECISION_REJECTED,
					"decision_reason": constant.DIGISIGN_VERIFICATION_FAILED,
					"journey":         "DGS",
					"callback":        false,
				},
				"server_time": utils.GenerateTimeNow(),
			})

			var worker []entity.TrxWorker

			if dataWorker.TransactionType == "USE_LIMIT" {
				paramRTL, _ := json.Marshal(map[string]interface{}{
					"amount": dataWorker.AF, "customer_id": dataWorker.CustomerID, "source_process": "LOS", "prospect_id": prospectID,
				})

				worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_UNPROCESS, EndPointTarget: os.Getenv("KREDITMU_RETURN_URL"),
					EndPointMethod: "PUT", Payload: string(paramRTL),
					ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAE, MaxRetry: 6, CountRetry: 0,
					Category: constant.WORKER_CATEGORY_KREDITMU, Action: constant.RETURN_LIMIT, Sequence: 1,
				})

			}

			var (
				sequence       int = 1
				activityWorker     = constant.WORKER_UNPROCESS
			)

			if len(worker) > 0 {
				sequence = 2
				activityWorker = constant.WORKER_IDLE
			}

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: activityWorker, EndPointTarget: os.Getenv("CALLBACK_URL_DIGISIGN"),
				EndPointMethod: constant.METHOD_POST, Payload: string(param), Header: string(headerParam),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
				Category: constant.WORKER_CATEGORY_DIGISIGN, Action: constant.CALLBACK_STATUS_1200, Sequence: sequence,
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_IDLE, EndPointTarget: dataWorker.CallbackUrl,
				EndPointMethod: constant.METHOD_POST, Payload: string(paramCallbackStatus), Header: string(headerParam),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
				Category: constant.CATEGORY_FRONT_APP, Action: constant.ACTION_CALLBACK, Sequence: sequence + 1,
			})

			go u.repository.SaveToWorker(worker)

		}

	}

	return

}

func RegisterMappingResponse(data response.RegisterResponse, prospectID string) (resp response.DataRegisterResponse) {

	resp = response.DataRegisterResponse{
		ProspectID: prospectID,
		Decision:   constant.DECISION_REJECT,
	}

	switch data.JsonFile.Result {

	case constant.CODE_00:

		resp.Decision = constant.DECISION_PASS

		if data.JsonFile.Notif == constant.REASON_SUCCESS_REGISTER {
			resp.DecisionReason = constant.REASON_SUCCESS_REGISTER
			resp.Code = constant.CODE_SUCCESS_REGISTER
			return

		} else if data.JsonFile.Notif == constant.REASON_EMAIL_REGISTERED {
			resp.DecisionReason = constant.REASON_EMAIL_REGISTERED
			resp.Code = constant.CODE_EMAIL_REGISTERED
			return

		} else if data.JsonFile.Notif == constant.REASON_REGISTERED_BEFORE {
			resp.DecisionReason = constant.REASON_REGISTERED_BEFORE
			resp.Code = constant.CODE_REGISTERED_BEFORE
			resp.IsRegistered = true
			return

		} else {
			resp.DecisionReason = constant.REASON_REGISTERED
			resp.Code = constant.CODE_REGISTERED
			resp.IsRegistered = true
			return
		}

	case constant.CODE_05:
		if data.JsonFile.Notif == constant.REASON_REGISTER_EXIST {
			resp.DecisionReason = constant.REASON_REGISTER_EXIST
			resp.Code = constant.CODE_REGISTER_EXIST
			resp.IsRegistered = true
			resp.Decision = constant.DECISION_PASS

			return
		}

	case constant.CODE_12:

		resp.DecisionReason = fmt.Sprintf("%s %s", data.JsonFile.Notif, data.JsonFile.Info)

		if data.JsonFile.Info == constant.REASON_REGISTER_FAILED {
			resp.Code = constant.CODE_REGISTER_FAILED
			return

		} else if strings.ReplaceAll(data.JsonFile.Notif, `\`, "") == constant.REASON_REGISTER_FAILED_NIK {
			resp.Code = constant.CODE_REGISTER_FAILED_NIK
			resp.DecisionReason = constant.REASON_REGISTER_FAILED_NIK
			return

		} else if data.JsonFile.Info == constant.REASON_REGISTER_FAILED_NOFACE_SELFIE {
			resp.Code = constant.CODE_REGISTER_FAILED_NOFACE_SELFIE
			return

		} else if data.JsonFile.Info == constant.REASON_REGISTER_FAILED_NOFACE_KTP {
			resp.Code = constant.CODE_REGISTER_FAILED_NOFACE_KTP
			return
		} else if data.JsonFile.Info == constant.REASON_REGISTER_MIN_LIGHT {
			resp.Code = constant.CODE_REGISTER_MIN_LIGHT
			return

		} else if data.JsonFile.Info == constant.REASON_REGISTER_FAILED_FACE_INVALID {
			resp.Code = constant.CODE_REGISTER_FAILED_FACE_INVALID
			return

		} else if data.JsonFile.Info == constant.REASON_REGISTER_FAILED_MIN_SIZE {
			resp.Code = constant.CODE_REGISTER_FAILED_MIN_SIZE
			return

		} else if data.JsonFile.Info == constant.REASON_REGISTER_FAILED_JPEG_FORMAT {
			resp.Code = constant.CODE_REGISTER_FAILED_JPEG_FORMAT
			return

		} else {
			resp.Code = constant.CODE_REGISTER_FAILED_NAMA
			resp.DecisionReason = constant.REASON_REGISTER_FAILED_NAMA
			return
		}

	case constant.CODE_14:

		if strings.Contains(data.JsonFile.Notif, "NIK sudah terdaftar") {
			resp.Code = constant.CODE_REGISTER_FAILED_NIK_REGISTERED
			resp.DecisionReason = fmt.Sprintf("%s %s", constant.REASON_REGISTER_FAILED_NIK_REGISTERED, data.JsonFile.EmailRegistered)
			return

		} else {
			resp.Code = constant.CODE_REGISTER_FAILED_EMAIL_REGISTERED
			resp.DecisionReason = constant.REASON_REGISTER_FAILED_EMAIL_REGISTERED
			return
		}

	case constant.CODE_15:

		resp.Code = constant.CODE_REGISTER_FAILED_MOBILE_PHONE_REGISTERED
		resp.DecisionReason = constant.REASON_REGISTER_FAILED_MOBILE_PHONE_REGISTERED
		return

	case constant.CODE_20:

		resp.Code = constant.CODE_REGISTER_FAILED_MAX_LIMIT
		resp.DecisionReason = constant.REASON_REGISTER_FAILED_MAX_LIMIT
		return

	case constant.CODE_91:

		resp.Code = constant.CODE_REGISTER_FAILED_TIMEOUT
		resp.DecisionReason = constant.REASON_REGISTER_FAILED_TIMEOUT
		return

	default:

		resp.Code = constant.CODE_REGISTER_GENERAL_ERROR
		resp.DecisionReason = constant.REASON_REGISTER_GENERAL_ERROR
		return
	}

	return

}

func (u packages) GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpByte, selfieByte, signatureByte, npwpByte []byte, err error) {

	var (
		ktpBase64       string
		selfieBase64    string
		signatureBase64 string
		npwpBase64      string
	)

	ktpMedia := GetIsMedia(ktpUrl)

	if ktpMedia {
		ktpBase64, err = u.usecase.DecodeMedia(ktpUrl, prospectID)
		if err != nil {
			return
		}

	} else {
		ktpBase64, err = DecodeNonMedia(ktpUrl)
		if err != nil {
			return
		}
	}

	ktpByte, err = utils.Base64Decode(ktpBase64)

	if err != nil {
		return
	}

	selfieMedia := GetIsMedia(selfieUrl)

	if selfieMedia {
		selfieBase64, err = u.usecase.DecodeMedia(selfieUrl, prospectID)
		if err != nil {
			return
		}

	} else {
		selfieBase64, err = DecodeNonMedia(selfieUrl)
		if err != nil {
			return
		}
	}

	selfieByte, err = utils.Base64Decode(selfieBase64)

	if err != nil {
		return
	}

	if signatureBase64 != "" {
		signatureMedia := GetIsMedia(signatureUrl)

		if signatureMedia {
			signatureBase64, err = u.usecase.DecodeMedia(signatureUrl, prospectID)
			if err != nil {
				return
			}

		} else {
			signatureBase64, err = DecodeNonMedia(signatureUrl)
			if err != nil {
				return
			}
		}

		signatureByte, err = utils.Base64Decode(signatureBase64)

		if err != nil {
			return
		}

	}

	if npwpUrl != "" {
		photoNpwp := GetIsMedia(npwpUrl)

		if photoNpwp {
			npwpBase64, err = u.usecase.DecodeMedia(npwpUrl, prospectID)
			if err != nil {
				return
			}

		} else {
			npwpBase64, err = DecodeNonMedia(npwpUrl)
			if err != nil {
				return
			}
		}

		npwpByte, err = utils.Base64Decode(npwpBase64)

		if err != nil {
			return
		}

	}

	return
}

func (u multiUsecase) Activation(req request.ActivationRequest) (data response.DataActivationResponse, err error) {

	jsonFile, _ := json.Marshal(map[string]interface{}{
		"JSONFile": map[string]interface{}{
			"userid":     os.Getenv("DIGISIGN_USER_ID"),
			"email_user": req.Email,
		},
	})

	param := map[string]string{
		"jsonfield": string(jsonFile),
	}

	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": "Bearer " + os.Getenv("DIGISIGN_TOKEN"),
	}

	var digiResp response.ActivationResponse

	dummyActive, _ := strconv.ParseBool(os.Getenv("DUMMY_ACTIVE"))

	if dummyActive {

		var dummy entity.DigisignDummy

		dummy, err = u.repository.GetDigisignDummy(req.Email, "ACTIVATION")

		json.Unmarshal([]byte(dummy.Response), &digiResp)

		go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
			ProspectID: req.ProspectID,
			Response:   dummy.Response,
			Activity:   "ACTIVATION",
			Link:       digiResp.JsonFile.Link,
		})

		data = ActivationMappingResponse(digiResp, req.ProspectID)

		return
	}

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	resp, err := u.httpclient.ActivationAPI(os.Getenv("ACTIVATION_URL"), constant.METHOD_POST, param, header, timeOut, req.ProspectID)

	if err != nil {
		return
	}

	json.Unmarshal(resp.Body(), &digiResp)

	go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
		ProspectID: req.ProspectID,
		Response:   string(resp.Body()),
		Activity:   "ACTIVATION",
		Link:       digiResp.JsonFile.Link,
	})

	isError := false

	if resp.StatusCode() != 200 {
		isError = true
	}

	logs := map[string]interface{}{
		"ID":            req.ProspectID,
		"request":       param,
		"response":      digiResp,
		"response_code": resp.StatusCode(),
		"url":           os.Getenv("ACTIVATION_URL"),
		"response_time": fmt.Sprintf("%dms", resp.Time().Milliseconds()),
	}

	go config.SetCustomLog("API_DIGISIGN", isError, logs, "ACTIVATION-API")

	data = ActivationMappingResponse(digiResp, req.ProspectID)

	go u.usecase.CallbackDigisignActivation(data, req.ProspectID)

	return
}

func (u usecase) CallbackDigisignActivation(data response.DataActivationResponse, prospectID string) (err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	status := constant.WAITING_ACTIVATION
	activity := constant.ACTIVITY_ONPROSES
	action := constant.CALLBACK_STATUS_1202

	if data.Decision == constant.DECISION_REJECT {
		status = constant.REGISTER_FAILED
		activity = constant.ACTIVITY_FINISHED
		action = constant.CALLBACK_STATUS_1201
	}

	dataWorker, _ := u.repository.GetDataWorker(prospectID)

	if dataWorker.TransactionType != "" {
		header := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": os.Getenv("AUTH_KPM"),
		}

		headerParam, _ := json.Marshal(header)

		param, _ := json.Marshal(map[string]interface{}{
			"prospect_id": prospectID,
			"code":        data.Code,
			"status":      status,
			"activity":    activity,
		})

		var worker []entity.TrxWorker

		if data.Decision == constant.DECISION_REJECT && dataWorker.TransactionType == "USE_LIMIT" {
			paramRTL, _ := json.Marshal(map[string]interface{}{
				"amount": dataWorker.AF, "customer_id": dataWorker.CustomerID, "source_process": "LOS", "prospect_id": prospectID,
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_UNPROCESS, EndPointTarget: os.Getenv("KREDITMU_RETURN_URL"),
				EndPointMethod: "PUT", Payload: string(paramRTL),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAE, MaxRetry: 6, CountRetry: 0,
				Category: constant.WORKER_CATEGORY_KREDITMU, Action: constant.RETURN_LIMIT, Sequence: 1,
			})
		}

		var (
			sequence       interface{}
			activityWorker = constant.WORKER_UNPROCESS
		)

		if len(worker) > 0 {
			sequence = 2
			activityWorker = constant.WORKER_IDLE
		}

		worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: activityWorker, EndPointTarget: os.Getenv("CALLBACK_URL_DIGISIGN"),
			EndPointMethod: constant.METHOD_POST, Payload: string(param), Header: string(headerParam),
			ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
			Category: constant.WORKER_CATEGORY_DIGISIGN, Action: action, Sequence: sequence,
		})

		if data.Decision == constant.DECISION_REJECT {

			var newSequence = 1
			if sequence != nil {
				newSequence = sequence.(int) + 1
			}

			paramCallbackStatus, _ := json.Marshal(map[string]interface{}{
				"messages": "LOS DigiSign",
				"errors":   nil,
				"data": map[string]interface{}{
					"prospect_id":     prospectID,
					"code":            data.Code,
					"decision":        constant.DECISION_REJECTED,
					"decision_reason": constant.DIGISIGN_VERIFICATION_FAILED,
					"journey":         "DGS",
					"callback":        false,
				},
				"server_time": utils.GenerateTimeNow(),
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_IDLE, EndPointTarget: dataWorker.CallbackUrl,
				EndPointMethod: constant.METHOD_POST, Payload: string(paramCallbackStatus), Header: string(headerParam),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
				Category: constant.CATEGORY_FRONT_APP, Action: constant.ACTION_CALLBACK, Sequence: newSequence,
			})

		}

		go u.repository.SaveToWorker(worker)

	}

	return

}

func ActivationMappingResponse(data response.ActivationResponse, prospectID string) (resp response.DataActivationResponse) {

	resp = response.DataActivationResponse{
		ProspectID: prospectID,
		Decision:   constant.DECISION_PASS,
	}

	switch data.JsonFile.Result {

	case constant.CODE_00:
		resp.DecisionReason = constant.REASON_SUCCESS_ACTIVATION
		resp.Code = constant.CODE_SUCCESS_ACTIVATION
		resp.Link = data.JsonFile.Link
		return

	case constant.CODE_06:
		resp.Decision = constant.DECISION_REJECT
		resp.DecisionReason = constant.REASON_ACTIVATION_FAILED_GENERAL_ERROR
		resp.Code = constant.CODE_ACTIVATION_FAILED_GENERAL_ERROR
		return

	case constant.CODE_14:
		resp.DecisionReason = constant.REASON_ACTIVATION_EXIST
		resp.Code = constant.CODE_ACTIVATION_EXIST
		resp.IsRegistered = true
		return

	default:
		resp.Decision = constant.DECISION_REJECT
		resp.DecisionReason = constant.REASON_ACTIVATION_GENERAL_ERROR
		resp.Code = constant.CODE_ACTIVATION_GENERAL_ERROR
		return
	}

}

func (u multiUsecase) ActivationRedirect(msg string) (data response.DataSignDocResponse, err error) {

	decodeValue, _ := base64.StdEncoding.DecodeString(msg)

	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(os.Getenv("DIGISIGN_AES_KEY")))

	var activationCallback response.ActivationCallbackResponse

	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(byteDecrypt, &activationCallback)

	logs := map[string]interface{}{
		"response": activationCallback,
		"msg":      msg,
	}

	go config.SetCustomLog("API_DIGISIGN", false, logs, "ACTIVATION-MSG")

	if activationCallback.Result == constant.CODE_00 && activationCallback.Notif == constant.ACTIVATION_CALLLBACK_SUCCESS {

		var (
			dataCustomer entity.CallbackData
			sendDoc      response.DataSendDocResponse
			signDoc      response.DataSignDocResponse
		)

		dataCustomer, err = u.repository.GetCustomerPersonalByEmailAndNik(activationCallback.Email, activationCallback.NIK)
		if err != nil {
			data = response.DataSignDocResponse{
				ProspectID:     dataCustomer.ProspectID,
				Code:           constant.CODE_EXPIRED_ACTIVATION,
				DecisionReason: constant.DECISION_EXPIRED_ACTIVATION,
				Decision:       constant.DECISION_REJECT,
				Link:           dataCustomer.RedirectFailedUrl,
			}
			err = fmt.Errorf(constant.EXPIRED)
			return
		}

		if dataCustomer.DiffTime > 5 && dataCustomer.DiffTime <= 10 {
			data = response.DataSignDocResponse{
				ProspectID:     dataCustomer.ProspectID,
				Code:           constant.CODE_EXPIRED_ACTIVATION,
				DecisionReason: constant.DECISION_EXPIRED_ACTIVATION,
				Decision:       constant.DECISION_REJECT,
				Link:           dataCustomer.RedirectFailedUrl,
			}

			return

		} else if dataCustomer.DiffTime > 10 {

			data = response.DataSignDocResponse{
				ProspectID:     dataCustomer.ProspectID,
				Code:           constant.CODE_EXPIRED_ACTIVATION,
				DecisionReason: constant.DECISION_EXPIRED_ACTIVATION,
				Decision:       constant.DECISION_REJECT,
				Link:           dataCustomer.RedirectFailedUrl,
			}

			err = fmt.Errorf(constant.EXPIRED)

			return
		}

		logs := map[string]interface{}{
			"ID":       dataCustomer.ProspectID,
			"response": activationCallback,
			"msg":      msg,
		}

		go config.SetCustomLog("API_DIGISIGN", false, logs, "ACTIVATION-CALLBACK-API")

		cache := utils.GetCache()

		keyCache := dataCustomer.ProspectID + "ACT-Redirect"

		value, _ := cache.Get(keyCache)

		if value != nil {
			jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(value, &data)
			return
		}

		var details []entity.TrxDetail

		go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
			ProspectID: dataCustomer.ProspectID,
			Response:   string(byteDecrypt),
			Activity:   "CALLBACK_ACTIVATION",
		})

		details = append(details, entity.TrxDetail{
			ProspectID: dataCustomer.ProspectID, StatusProcess: "ONP", Activity: "PRCD", Decision: "PAS",
			RuleCode: "4225", SourceDecision: "ACT", NextStep: nil, CreatedBy: "SYSTEM",
		})

		err = u.repository.UpdateStatusDigisignActivation(activationCallback.Email, activationCallback.NIK, dataCustomer.ProspectID, details)

		if err != nil {
			return
		}

		details = []entity.TrxDetail{}

		// send doc
		sendDoc, err = u.packages.SendDoc(request.SendDoc{
			ProspectID: dataCustomer.ProspectID,
			Email:      activationCallback.Email,
			IDKtp:      activationCallback.NIK,
		})

		if err != nil {
			return
		}

		info, _ := json.Marshal(response.SendDocInfo{
			DocumentID:  sendDoc.DocumentID,
			AgreementNo: sendDoc.AgreementNo,
		})

		if sendDoc.Decision == constant.DECISION_REJECT {
			details = append(details, entity.TrxDetail{
				ProspectID: dataCustomer.ProspectID, StatusProcess: "FIN", Activity: "STOP", Decision: "REJ",
				RuleCode: sendDoc.Code, SourceDecision: "SND", NextStep: nil, CreatedBy: "SYSTEM",
			})

			data = response.DataSignDocResponse{
				ProspectID:     sendDoc.ProspectID,
				Code:           sendDoc.Code,
				Decision:       sendDoc.Decision,
				DecisionReason: sendDoc.DecisionReason,
				Link:           dataCustomer.RedirectFailedUrl,
			}

			cacheData, _ := json.Marshal(data)

			cache.Set(keyCache, cacheData)

			err = u.repository.SaveTrx(details)

			return
		}

		details = append(details, entity.TrxDetail{
			ProspectID: dataCustomer.ProspectID, StatusProcess: "ONP", Activity: "PRCD", Decision: "PAS",
			RuleCode: sendDoc.Code, SourceDecision: "SND", NextStep: "SID", CreatedBy: "SYSTEM", Info: string(info),
		})

		// sign doc
		signDoc, err = u.packages.SignDocument(request.JsonFileSign{
			UserID:     os.Getenv("DIGISIGN_USER_ID"),
			DocumentID: sendDoc.DocumentID,
			Email:      activationCallback.Email,
			ViewOnly:   false,
		}, dataCustomer.ProspectID)

		if err != nil {
			return
		}

		data = signDoc

		if signDoc.Decision == constant.DECISION_REJECT {

			details = append(details, entity.TrxDetail{
				ProspectID: dataCustomer.ProspectID, StatusProcess: "FIN", Activity: "STOP", Decision: "REJ",
				RuleCode: signDoc.Code, SourceDecision: "SID", NextStep: nil, CreatedBy: "SYSTEM",
			})

			data = response.DataSignDocResponse{
				ProspectID:     signDoc.ProspectID,
				Code:           signDoc.Code,
				Decision:       signDoc.Decision,
				DecisionReason: signDoc.DecisionReason,
				Link:           dataCustomer.RedirectFailedUrl,
			}

			cacheData, _ := json.Marshal(data)

			cache.Set(keyCache, cacheData)

			err = u.repository.SaveTrx(details)

			return

		}
		// check if SND exist
		sndExist := u.repository.CheckSND(dataCustomer.ProspectID)
		if sndExist <= 0 {
			details = append(details, entity.TrxDetail{
				ProspectID: dataCustomer.ProspectID, StatusProcess: "ONP", Activity: "UNPR", Decision: "CPR",
				RuleCode: signDoc.Code, SourceDecision: "SID", NextStep: nil, CreatedBy: "SYSTEM", Info: sendDoc.DocumentID + ".pdf",
			})
		}else {
			details = []entity.TrxDetail{}
		}


		cacheData, _ := json.Marshal(data)

		cache.Set(keyCache, cacheData)

		err = u.repository.SaveTrx(details)
		if err != nil {
			return
		}

		return

	}

	return
}

func (u usecase) GeneratePK(prospectID string) (document []byte, docID string, agreementNo string, err error) {

	param, _ := json.Marshal(map[string]interface{}{
		"ProspectID": prospectID,
	})

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	resp, _ := u.httpclient.DocumentAPI(os.Getenv("GENERATE_PK_URL"), constant.METHOD_POST, param, map[string]string{}, timeOut, prospectID)

	logs := map[string]interface{}{
		"ID":            prospectID,
		"request":       param,
		"response_code": resp.StatusCode(),
		"url":           os.Getenv("GENERATE_PK_URL"),
		"response_time": fmt.Sprintf("%dms", resp.Time().Milliseconds()),
	}

	if resp.StatusCode() != 200 {
		go config.SetCustomLog("API_DIGISIGN", true, logs, "GENERATE-PK-API")
		err = fmt.Errorf("failed Generate PK")
		return
	}

	var documentData response.DocumentGenerateResponse

	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(jsoniter.Get(resp.Body(), "data").ToString()), &documentData)

	document, err = utils.Base64Decode(documentData.DocumentPK)
	if err != nil {
		return
	}

	docID = documentData.DocumentID
	agreementNo = documentData.AgreementNo

	logs["document_Id"] = docID
	logs["agreement_no"] = agreementNo

	go config.SetCustomLog("API_DIGISIGN", false, logs, "GENERATE-PK-API")

	return
}

func (u packages) SendDoc(req request.SendDoc) (data response.DataSendDocResponse, err error) {

	sendDoc, err := u.repository.GetSendDocData(req.ProspectID)
	if err != nil {
		return
	}

	document, documentID, agreementNo, err := u.usecase.GeneratePK(req.ProspectID)
	if err != nil {
		return
	}

	dataFile := request.DataFile{
		DocumentPK: document,
		DocumentID: documentID,
	}

	jsonFile := request.SendDocRequest{
		UserID:         os.Getenv("DIGISIGN_USER_ID"),
		DocumentID:     documentID,
		Payment:        os.Getenv("PAYMENT_METHOD"),
		Redirect:       true,
		Branch:         sendDoc.BranchID,
		SequenceOption: false,
		SendTo: []request.SendTo{
			{Name: sendDoc.LegalName, Email: req.Email},
			{Name: sendDoc.NameBM, Email: sendDoc.EmailBM},
		},
		ReqSign: []request.ReqSign{
			{
				Name: sendDoc.NameBM, Email: sendDoc.EmailBM, AksiTtd: "at", Kuser: sendDoc.Kuser, User: "prf1",
				Page: "1", Llx: "323", Lly: "135", Urx: "420", Ury: "184", Visible: "1",
			},
			{
				Name: sendDoc.LegalName, Email: req.Email, AksiTtd: "mt", User: "ttd1", Page: "1", Llx: "458",
				Lly: "135", Urx: "557", Ury: "184", Visible: "1",
			},
			{
				Name: sendDoc.LegalName, Email: req.Email, AksiTtd: "mt", User: "ttd2", Page: "5", Llx: "70",
				Lly: "356.7", Urx: "183", Ury: "457.5", Visible: "1",
			},
			{
				Name: sendDoc.NameBM, Email: sendDoc.EmailBM, AksiTtd: "at", Kuser: sendDoc.Kuser, User: "prf2",
				Page: "5", Llx: "428.4", Lly: "356.7", Urx: "541.4", Ury: "457.5", Visible: "1",
			},
			{
				Name: sendDoc.NameBM, Email: sendDoc.EmailBM, AksiTtd: "at", Kuser: sendDoc.Kuser, User: "prf3",
				Page: "7", Llx: "33", Lly: "448.6", Urx: "126.7", Ury: "495.4", Visible: "1",
			},
			{
				Name: sendDoc.LegalName, Email: req.Email, AksiTtd: "mt", User: "ttd3", Page: "7",
				Llx: "457", Lly: "448.6", Urx: "580", Ury: "495.4", Visible: "1",
			},
			{
				Name: sendDoc.LegalName, Email: req.Email, AksiTtd: "mt", User: "ttd4", Page: "9",
				Llx: "71.3", Lly: "251", Urx: "130", Ury: "283", Visible: "1",
			},
			{
				Name: sendDoc.NameBM, Email: sendDoc.EmailBM, AksiTtd: "at", Kuser: sendDoc.Kuser, User: "prf4",
				Page: "9", Llx: "445", Lly: "251", Urx: "546", Ury: "283", Visible: "1",
			},
			{
				Name: sendDoc.LegalName, Email: req.Email, AksiTtd: "mt", User: "ttd5", Page: "10",
				Llx: "31", Lly: "180", Urx: "118", Ury: "276.5", Visible: "1",
			},
		},
		SigningSeq: 0,
	}

	jsonField, _ := json.Marshal(request.JsonFile{
		JsonFile: jsonFile,
	})

	param := map[string]string{
		"jsonfield": string(jsonField),
	}

	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": "Bearer " + os.Getenv("DIGISIGN_TOKEN"),
	}

	var digiResponse response.SendDocResponse

	dummyActive, _ := strconv.ParseBool(os.Getenv("DUMMY_ACTIVE"))

	if dummyActive {

		var dummy entity.DigisignDummy

		dummy, err = u.repository.GetDigisignDummy(req.Email, "SEND_DOC")

		json.Unmarshal([]byte(dummy.Response), &digiResponse)

		data = SendDocMappingResponse(digiResponse, req.ProspectID)

		go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
			ProspectID: req.ProspectID,
			Response:   dummy.Response,
			Activity:   "SEND_DOC",
		})

		data.DocumentID = documentID
		data.AgreementNo = agreementNo

		return
	}

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_95S"))

	resp, err := u.httpclient.SendDocAPI(os.Getenv("SEND_DOC_URL"), constant.METHOD_POST, param, header, timeOut, dataFile, req.ProspectID)

	logs := map[string]interface{}{
		"ID":            req.ProspectID,
		"request":       param,
		"response":      string(resp.Body()),
		"response_code": resp.StatusCode(),
		"url":           os.Getenv("SEND_DOC_URL"),
		"response_time": fmt.Sprintf("%dms", resp.Time().Milliseconds()),
	}

	if err != nil {
		go config.SetCustomLog("API_DIGISIGN", true, logs, "SEND_DOC-API")
		return
	}

	json.Unmarshal(resp.Body(), &digiResponse)

	go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
		ProspectID: req.ProspectID,
		Response:   string(resp.Body()),
		Activity:   "SEND_DOC",
	})

	logs["document_Id"] = documentID
	logs["agreement_no"] = agreementNo
	logs["response"] = digiResponse

	go config.SetCustomLog("API_DIGISIGN", false, logs, "SEND_DOC-API")

	data = SendDocMappingResponse(digiResponse, req.ProspectID)

	data.DocumentID = documentID
	data.AgreementNo = agreementNo

	go u.usecase.CallbackDigisignSendDoc(data, req.ProspectID)

	return
}

func (u usecase) CallbackDigisignSendDoc(data response.DataSendDocResponse, prospectID string) (err error) {

	if data.Decision == constant.DECISION_REJECT {

		dataWorker, _ := u.repository.GetDataWorker(prospectID)

		if dataWorker.TransactionType != "" {
			timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

			header := map[string]string{
				"Content-Type":  "application/json",
				"Authorization": os.Getenv("AUTH_KPM"),
			}

			headerParam, _ := json.Marshal(header)

			param, _ := json.Marshal(map[string]interface{}{
				"prospect_id": prospectID,
				"code":        data.Code,
				"status":      constant.SEND_DOC_FAILED,
				"activity":    constant.ACTIVITY_FINISHED,
			})

			paramCallbackStatus, _ := json.Marshal(map[string]interface{}{
				"messages": "LOS DigiSign",
				"errors":   nil,
				"data": map[string]interface{}{
					"prospect_id":     prospectID,
					"code":            data.Code,
					"decision":        constant.DECISION_REJECTED,
					"decision_reason": constant.SEND_DOC_FAILED,
					"journey":         "DGS",
					"callback":        false,
				},
				"server_time": utils.GenerateTimeNow(),
			})

			var worker []entity.TrxWorker

			if dataWorker.TransactionType == "USE_LIMIT" {
				paramRTL, _ := json.Marshal(map[string]interface{}{
					"amount": dataWorker.AF, "customer_id": dataWorker.CustomerID, "source_process": "LOS", "prospect_id": prospectID,
				})

				worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_UNPROCESS, EndPointTarget: os.Getenv("KREDITMU_RETURN_URL"),
					EndPointMethod: "PUT", Payload: string(paramRTL),
					ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAE, MaxRetry: 6, CountRetry: 0,
					Category: constant.WORKER_CATEGORY_KREDITMU, Action: constant.RETURN_LIMIT, Sequence: 1,
				})

			}

			var (
				sequence       int = 1
				activityWorker     = constant.WORKER_UNPROCESS
			)

			if len(worker) > 0 {
				sequence = 2
				activityWorker = constant.WORKER_IDLE
			}

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: activityWorker, EndPointTarget: os.Getenv("CALLBACK_URL_DIGISIGN"),
				EndPointMethod: constant.METHOD_POST, Payload: string(param), Header: string(headerParam),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
				Category: constant.WORKER_CATEGORY_DIGISIGN, Action: constant.CALLBACK_STATUS_1210, Sequence: sequence,
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_IDLE, EndPointTarget: dataWorker.CallbackUrl,
				EndPointMethod: constant.METHOD_POST, Payload: string(paramCallbackStatus), Header: string(headerParam),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
				Category: constant.CATEGORY_FRONT_APP, Action: constant.ACTION_CALLBACK, Sequence: sequence + 1,
			})

			go u.repository.SaveToWorker(worker)

		}

	}

	return

}

func SendDocMappingResponse(data response.SendDocResponse, prospectID string) (resp response.DataSendDocResponse) {

	resp = response.DataSendDocResponse{
		ProspectID: prospectID,
		Decision:   constant.DECISION_REJECT,
	}

	switch data.JsonFile.Result {

	case constant.CODE_00:
		if data.JsonFile.Notif == constant.REASON_SUCCESS_SEND_DOC && data.JsonFile.Info == "" {
			resp.Decision = constant.DECISION_PASS
			resp.DecisionReason = constant.REASON_SUCCESS_SEND_DOC
			resp.Code = constant.CODE_SUCCESS_SEND_DOC
			return

		} else {
			resp.Decision = constant.DECISION_PASS
			resp.DecisionReason = fmt.Sprintf("%s %s", data.JsonFile.Notif, data.JsonFile.Info)
			resp.Code = constant.CODE_SEND_DOC_SUCCESS_WITH_CONDITION
			return
		}

	case constant.CODE_05:

		if data.JsonFile.Notif == constant.REASON_SEND_DOC_FAILED_DOCID_NULL {
			resp.DecisionReason = constant.REASON_SEND_DOC_FAILED_DOCID_NULL
			resp.Code = constant.CODE_SEND_DOC_FAILED_DOCID_NULL
			return

		} else {
			resp.DecisionReason = constant.REASON_SEND_DOC_GENERAL_ERROR
			resp.Code = constant.CODE_SEND_DOC_GENERAL_ERROR
			return
		}

	case constant.CODE_15:

		resp.DecisionReason = fmt.Sprintf("%s %s", data.JsonFile.Notif, data.JsonFile.Info)

		if strings.Contains(data.JsonFile.Info, constant.FLAG_SEND_DOC_REREGISTRATION) {
			resp.Code = constant.CODE_SEND_DOC_REREGISTRATION
			return
		} else if strings.Contains(data.JsonFile.Info, constant.FLAG_SEND_DOC_LOGIN_WEB) {
			resp.Code = constant.CODE_SEND_DOC_LOGIN_WEB
			return
		} else {
			resp.Code = constant.CODE_SEND_DOC_UNREGISTERED
			return
		}

	case constant.CODE_17:
		resp.DecisionReason = constant.REASON_SEND_DOC_FAILED_DOCID_EXIST
		resp.Code = constant.CODE_SEND_DOC_FAILED_DOCID_EXIST
		return

	case constant.CODE_91:
		resp.DecisionReason = constant.REASON_SEND_DOC_TIMEOUT
		resp.Code = constant.CODE_SEND_DOC_TIMEOUT
		return

	default:
		resp.DecisionReason = constant.REASON_SEND_DOC_GENERAL_ERROR
		resp.Code = constant.CODE_SEND_DOC_GENERAL_ERROR
		return
	}

}

func (u packages) SignDocument(req request.JsonFileSign, prospectID string) (data response.DataSignDocResponse, err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	jsonField, _ := json.Marshal(request.JsonFile{
		JsonFile: req,
	})

	param := map[string]string{
		"jsonfield": string(jsonField),
	}

	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": "Bearer " + os.Getenv("DIGISIGN_TOKEN"),
	}

	var digiResp response.SignDocResponse

	dummyActive, _ := strconv.ParseBool(os.Getenv("DUMMY_ACTIVE"))

	if dummyActive {

		var dummy entity.DigisignDummy

		dummy, err = u.repository.GetDigisignDummy(req.Email, "SIGN_DOC")

		json.Unmarshal([]byte(dummy.Response), &digiResp)

		go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
			ProspectID: prospectID,
			Response:   dummy.Response,
			Activity:   "SIGN_DOC",
		})

		data = SignDocumentMappingResponse(digiResp, prospectID)

		return
	}

	signDoc, err := u.httpclient.SignDocAPI(os.Getenv("SIGN_DOC_URL"), constant.METHOD_POST, param, header, timeOut, prospectID)

	logs := map[string]interface{}{
		"ID":            prospectID,
		"request":       param,
		"response":      string(signDoc.Body()),
		"response_code": signDoc.StatusCode(),
		"url":           os.Getenv("SIGN_DOC_URL"),
		"response_time": fmt.Sprintf("%dms", signDoc.Time().Milliseconds()),
	}

	if err != nil {
		go config.SetCustomLog("API_DIGISIGN", true, logs, "SIGN_DOC-API")
		return
	}

	json.Unmarshal(signDoc.Body(), &digiResp)

	go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
		ProspectID: prospectID,
		Response:   string(signDoc.Body()),
		Activity:   "SIGN_DOC",
		Link:       digiResp.JsonFile.Link,
	})

	logs["response"] = digiResp

	go config.SetCustomLog("API_DIGISIGN", false, logs, "SIGN_DOC-API")

	data = SignDocumentMappingResponse(digiResp, prospectID)

	go u.usecase.CallbackDigisignSignDoc(data, prospectID)

	return
}

func (u usecase) CallbackDigisignSignDoc(data response.DataSignDocResponse, prospectID string) (err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	status := constant.WAITING_SIGN_DOC
	activity := constant.ACTIVITY_ONPROSES
	action := constant.CALLBACK_STATUS_1206

	if data.Decision == constant.DECISION_REJECT {
		status = constant.SEND_DOC_FAILED
		activity = constant.ACTIVITY_FINISHED
		action = constant.CALLBACK_STATUS_1210
	}

	dataWorker, _ := u.repository.GetDataWorker(prospectID)

	if dataWorker.TransactionType != "" {
		header := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": os.Getenv("AUTH_KPM"),
		}

		headerParam, _ := json.Marshal(header)

		param, _ := json.Marshal(map[string]interface{}{
			"prospect_id": prospectID,
			"code":        data.Code,
			"status":      status,
			"activity":    activity,
		})

		var worker []entity.TrxWorker

		if data.Decision == constant.DECISION_REJECT && dataWorker.TransactionType == "USE_LIMIT" {
			paramRTL, _ := json.Marshal(map[string]interface{}{
				"amount": dataWorker.AF, "customer_id": dataWorker.CustomerID, "source_process": "LOS", "prospect_id": prospectID,
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_UNPROCESS, EndPointTarget: os.Getenv("KREDITMU_RETURN_URL"),
				EndPointMethod: "PUT", Payload: string(paramRTL),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAE, MaxRetry: 6, CountRetry: 0,
				Category: constant.WORKER_CATEGORY_KREDITMU, Action: constant.RETURN_LIMIT, Sequence: 1,
			})
		}

		var (
			sequence       interface{}
			activityWorker = constant.WORKER_UNPROCESS
		)

		if len(worker) > 0 {
			sequence = 2
			activityWorker = constant.WORKER_IDLE
		}

		worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: activityWorker, EndPointTarget: os.Getenv("CALLBACK_URL_DIGISIGN"),
			EndPointMethod: constant.METHOD_POST, Payload: string(param), Header: string(headerParam),
			ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
			Category: constant.WORKER_CATEGORY_DIGISIGN, Action: action, Sequence: sequence,
		})

		if data.Decision == constant.DECISION_REJECT {

			var newSequence = 1
			if sequence != nil {
				newSequence = sequence.(int) + 1
			}

			paramCallbackStatus, _ := json.Marshal(map[string]interface{}{
				"messages": "LOS DigiSign",
				"errors":   nil,
				"data": map[string]interface{}{
					"prospect_id":     prospectID,
					"code":            data.Code,
					"decision":        constant.DECISION_REJECTED,
					"decision_reason": constant.DIGISIGN_VERIFICATION_FAILED,
					"journey":         "DGS",
					"callback":        false,
				},
				"server_time": utils.GenerateTimeNow(),
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_IDLE, EndPointTarget: dataWorker.CallbackUrl,
				EndPointMethod: constant.METHOD_POST, Payload: string(paramCallbackStatus), Header: string(headerParam),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
				Category: constant.CATEGORY_FRONT_APP, Action: constant.ACTION_CALLBACK, Sequence: newSequence,
			})

		}

		go u.repository.SaveToWorker(worker)

	}

	return

}

func SignDocumentMappingResponse(data response.SignDocResponse, prospectID string) (resp response.DataSignDocResponse) {

	resp = response.DataSignDocResponse{
		ProspectID: prospectID,
		Decision:   constant.DECISION_REJECT,
	}

	switch data.JsonFile.Result {

	case constant.CODE_00:
		resp.Decision = constant.DECISION_PASS
		resp.DecisionReason = constant.REASON_SUCCESS_SIGN_DOC
		resp.Code = constant.CODE_SUCCESS_SIGN_DOC
		resp.Link = data.JsonFile.Link
		return

	case constant.CODE_05:

		if data.JsonFile.Notif == constant.REASON_SIGN_DOC_INVALID {
			resp.DecisionReason = constant.REASON_SIGN_DOC_INVALID
			resp.Code = constant.CODE_SIGN_DOC_INVALID
			return
		} else if data.JsonFile.Notif == constant.REASON_SIGN_DOC_EXIST {
			resp.DecisionReason = constant.REASON_SIGN_DOC_EXIST
			resp.Code = constant.CODE_SIGN_DOC_EXIST
			return
		} else if data.JsonFile.Notif == constant.REASON_SIGN_DOC_NOTFOUND {
			resp.DecisionReason = constant.REASON_SIGN_DOC_NOTFOUND
			resp.Code = constant.CODE_SIGN_DOC_NOTFOUND
			return
		} else {
			resp.Code = constant.CODE_SIGN_DOC_GENERAL_ERROR
			resp.DecisionReason = constant.REASON_SIGN_DOC_GENERAL_ERROR
			return
		}

	default:
		resp.Code = constant.CODE_SIGN_DOC_GENERAL_ERROR
		resp.DecisionReason = constant.REASON_SIGN_DOC_GENERAL_ERROR
		return
	}

}

func (u multiUsecase) SignCallback(msg string) (upload response.MediaServiceResponse, redirectUrl string, err error) {

	decodeValue, _ := base64.StdEncoding.DecodeString(msg)

	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(os.Getenv("DIGISIGN_AES_KEY")))

	var signCallback response.SignCallback

	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(byteDecrypt, &signCallback)

	logs := map[string]interface{}{
		"response": signCallback,
		"msg":      msg,
	}

	go config.SetCustomLog("API_DIGISIGN", false, logs, "SIGN-MSG")

	if signCallback.StatusDocument == constant.SIGN_DOC_COMPLETE && signCallback.Result == constant.CODE_00 {

		var data entity.CallbackData

		data, err = u.repository.GetCustomerPersonalByEmail(signCallback.DocumentID)

		redirectUrl = data.RedirectFailedUrl

		logs := map[string]interface{}{
			"ID":       data.ProspectID,
			"response": signCallback,
			"msg":      msg,
		}

		go config.SetCustomLog("API_DIGISIGN", false, logs, "SIGN-CALLBACK-API")

		if data.Decision == "APR" {
			redirectUrl = data.RedirectSuccessUrl
			return
		} else {

			if err != nil {
				upload = response.MediaServiceResponse{
					Message: "LOS Digisign",
					Code:    constant.CODE_EXPIRED_SIGN_DOC,
				}

				err = fmt.Errorf(constant.EXPIRED)

				return
			}

			if data.DiffTime > 15 && data.DiffTime <= 30 {
				upload = response.MediaServiceResponse{
					Message: "LOS Digisign",
					Code:    constant.CODE_EXPIRED_SIGN_DOC,
				}

				return
			} else if data.DiffTime > 30 {
				upload = response.MediaServiceResponse{
					Message: "LOS Digisign",
					Code:    constant.CODE_EXPIRED_SIGN_DOC,
				}
				return
			}

		}

		var download string

		download, err = u.usecase.DownloadDoc(data.ProspectID, request.DownloadRequest{
			DocumentID: signCallback.DocumentID,
			UserID:     os.Getenv("DIGISIGN_USER_ID"),
		})

		fmt.Println("DOWNLOAD DOC : " )

		if err != nil {
			return
		}

		upload, err = u.usecase.UploadDoc(data.ProspectID, download)
		fmt.Println("DOC PK : ", upload.Data)
		log.Print(" UPLOAD : ", upload.Data)
		log.Print(" Error Upload : ", upload.Errors)
		if err != nil {
			return
		}

		details, _ := u.repository.GetAgreementNo(data.ProspectID)

		var info response.SendDocInfo

		if details.Info != nil {
			json.Unmarshal([]byte(details.Info.(string)), &info)
		}

		doc := entity.TteDocPk{
			ID:          uuid.New().String(),
			ProspectID:  data.ProspectID,
			NoAgreement: info.AgreementNo,
			DocPKUrl:    upload.Data.MediaURL,
			Tipe:        "DIGISIGN",
			FilePath:    os.Getenv("SIGNED_PATH") + info.DocumentID + ".pdf",
		}

		fmt.Println("DOC PK : ", doc)
		log.Print(" DOC PK : ", doc)
		_ = u.repository.UpdateStatusDigisignSignDoc(entity.TrxDetail{
			ProspectID: data.ProspectID, StatusProcess: "FIN", Activity: "STOP", Decision: "APR", RuleCode: "4404",
			SourceDecision: "SID", CreatedBy: "SYSTEM", Info: signCallback.DocumentID + ".pdf",
		}, doc)

		redirectUrl = data.RedirectSuccessUrl

		go u.repository.SaveToTrxDigisign(entity.TrxDigisign{
			ProspectID: data.ProspectID,
			Response:   string(byteDecrypt),
			Activity:   "SIGN_CALLBACK",
		})

		go u.usecase.CallbackDigisignSignDocSuccess(data.ProspectID, upload.Data.MediaURL)

	}
	return
}

func (u usecase) CallbackDigisignSignDocSuccess(prospectID string, url string) (err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	dataWorker, _ := u.repository.GetDataWorker(prospectID)

	result := u.repository.CheckWorker1209(prospectID)

	if result > 0 {
		return
	}

	if dataWorker.TransactionType != "" {
		header := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": os.Getenv("AUTH_KPM"),
		}

		headerParam, _ := json.Marshal(header)

		param, _ := json.Marshal(map[string]interface{}{
			"prospect_id":      prospectID,
			"code":             constant.CODE_SIGNED,
			"status":           constant.SIGN_DOC_SUCCESS,
			"final_pk_tte_url": url,
			"activity":         constant.ACTIVITY_FINISHED,
		})

		paramEpo, _ := json.Marshal(map[string]interface{}{
			"prospect_id":     prospectID,
			"code":            constant.CODE_SIGNED,
			"decision":        "APPROVED",
			"decision_reason": constant.SIGN_DOC_SUCCESS,
		})

		paramStaging, _ := json.Marshal(map[string]interface{}{
			"prospect_id": prospectID, "target": "CONFINS", "source": "DGS",
		})

		var worker []entity.TrxWorker

		if dataWorker.TransactionType == "PRODUCT_LIMIT" {
			paramUSL, _ := json.Marshal(map[string]interface{}{
				"amount": dataWorker.AF, "application_source": "KPM", "customer_id": dataWorker.CustomerID, "source_process": "LOS", "tenor": dataWorker.TenorLimit, "prospect_id": prospectID,
			})

			worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_UNPROCESS, EndPointTarget: os.Getenv("KREDITMU_USE_URL"),
				EndPointMethod: "PUT", Payload: string(paramUSL),
				ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAE, MaxRetry: 6, CountRetry: 0,
				Category: constant.WORKER_CATEGORY_KREDITMU, Action: constant.USE_LIMIT, Sequence: 1,
			})

		}

		var (
			sequence       int = 1
			activityWorker     = constant.WORKER_UNPROCESS
		)
		if len(worker) > 0 {
			sequence = 2
			activityWorker = constant.WORKER_IDLE
		}

		worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: activityWorker, EndPointTarget: os.Getenv("CALLBACK_URL_DIGISIGN"),
			EndPointMethod: constant.METHOD_POST, Payload: string(param), Header: string(headerParam),
			ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
			Category: constant.WORKER_CATEGORY_DIGISIGN, Action: constant.CALLBACK_STATUS_1209, Sequence: sequence,
		})

		sequence = sequence + 1

		worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_IDLE, EndPointTarget: os.Getenv("EPO_URL"),
			EndPointMethod: constant.METHOD_POST, Payload: string(paramEpo),
			ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
			Category: constant.WORKER_CATEGORY_EPO, Action: constant.CREATE_EPO, Sequence: sequence,
		})

		sequence = sequence + 1

		worker = append(worker, entity.TrxWorker{ProspectID: prospectID, Activity: constant.WORKER_IDLE, EndPointTarget: os.Getenv("INSERT_STAGING_URL"),
			EndPointMethod: constant.METHOD_POST, Payload: string(paramStaging),
			ResponseTimeout: timeOut, APIType: constant.WORKER_TYPE_RAW, MaxRetry: 6, CountRetry: 0,
			Category: constant.WORKER_CATEGORY_CONFINS, Action: constant.INSERT_STAGING, Sequence: sequence,
		})

		go u.repository.SaveToWorker(worker)

	}

	return

}

func (u usecase) DownloadDoc(prospectID string, req request.DownloadRequest) (pdfBase64 string, err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	jsonFile, _ := json.Marshal(map[string]interface{}{
		"JSONFile": req,
	})

	param := map[string]string{
		"jsonfield": string(jsonFile),
	}

	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": "Bearer " + os.Getenv("DIGISIGN_TOKEN"),
	}

	restyResp, err := u.httpclient.DigiAPI(os.Getenv("DOWNLOAD_URL"), constant.METHOD_POST, param, "", header, timeOut, prospectID)
	var respDownload response.DownloadResponse
	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
		if err := json.Unmarshal(restyResp.Body(), &respDownload); err != nil {
			return pdfBase64, err
		}
	}

	logs := map[string]interface{}{
		"ID":            prospectID,
		"response_code": restyResp.StatusCode(),
		"url":           os.Getenv("DOWNLOAD_URL"),
		"response_time": fmt.Sprintf("%dms", restyResp.Time().Milliseconds()),
		"document_id":   req.DocumentID,
	}

	go config.SetCustomLog("API_DIGISIGN", false, logs, "DOWNLOAD-API")

	dec, err := base64.StdEncoding.DecodeString(respDownload.JsonFile.File)
	if err != nil {
		panic(err)
	}

	os.Remove(os.Getenv("SIGNED_PATH") + req.DocumentID + ".pdf")

	pdfBase64 = os.Getenv("SIGNED_PATH") + req.DocumentID + ".pdf"

	f, err := os.Create(pdfBase64)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
	return
}

func (u usecase) UploadDoc(prospectID string, fileName string) (uploadResp response.MediaServiceResponse, err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("DEFAULT_TIMEOUT_30S"))

	param := map[string]string{
		"type":         "perjanjian",
		"reference_no": prospectID,
	}
	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": os.Getenv("MEDIA_KEY"),
	}
	restyResp, err := u.httpclient.MediaClient(os.Getenv("UPLOAD_PLATFORM_URL"), http.MethodPost, param, fileName, header, timeOut, true, 3, prospectID)
	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
		if err := json.Unmarshal(restyResp.Body(), &uploadResp); err != nil {
			return uploadResp, err
		}
	}

	logs := map[string]interface{}{
		"ID":            prospectID,
		"response":      uploadResp,
		"response_code": restyResp.StatusCode(),
		"url":           os.Getenv("UPLOAD_PLATFORM_URL"),
		"response_time": fmt.Sprintf("%dms", restyResp.Time().Milliseconds()),
	}

	go config.SetCustomLog("API_DIGISIGN", false, logs, "UPLOAD-API")

	return
}

func DecodeNonMedia(url string) (base64Image string, err error) {

	image, err := http.Get(url)

	if err != nil {
		return
	}

	reader := bufio.NewReader(image.Body)
	ioutil, err := ioutil.ReadAll(reader)

	if err != nil {
		return
	}

	base64Image = base64.StdEncoding.EncodeToString(ioutil)

	return
}

func (u usecase) DecodeMedia(url string, customerID string) (base64Image string, err error) {

	timeOut, _ := strconv.Atoi(os.Getenv("MEDIA_TIMEOUT"))

	var decode response.ImageDecodeResponse

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": os.Getenv("MEDIA_KEY"),
	}

	image, err := u.httpclient.MediaClient(url+os.Getenv("MEDIA_PATH"), "GET", nil, "", header, timeOut, false, nil, customerID)

	logs := map[string]interface{}{
		"ID":            customerID,
		"response_code": image.StatusCode(),
		"url":           url + os.Getenv("UPLOAD_PLATFORM_URL"),
		"response_time": fmt.Sprintf("%dms", image.Time().Milliseconds()),
	}

	go config.SetCustomLog("API_DIGISIGN", false, logs, "MEDIA-API")

	if image.StatusCode() != 200 || err != nil {
		err = errors.New(constant.CONNECTION_ERROR)
		return
	}

	err = json.Unmarshal([]byte(image.Body()), &decode)

	if err != nil {
		err = fmt.Errorf("unmarshal error")
		return
	}

	base64Image = decode.Data.Encode

	return
}

func GetIsMedia(urlImage string) bool {

	urlMedia := strings.Split(os.Getenv("URL_MEDIA"), ",")

	for _, url := range urlMedia {
		if strings.Contains(urlImage, url) {
			return true
		}
	}

	return false
}

func (u usecase) DigisignCheck(email, prospectID string) (data response.DataDigisignCheck, err error) {

	getStatus, err := u.repository.GetTrxStatus(prospectID)

	if err != nil {
		err = fmt.Errorf(constant.DIGISIGN_RECORD_NOT_FOUND)
		return
	}

	if getStatus.SourceDecision == "ACT" || getStatus.SourceDecision == "RGS" || getStatus.SourceDecision == "SID" || getStatus.SourceDecision == "SND" {

		var step string
		var activationUrl, signDocUrl interface{}

		decision := "CREDIT_PROCESS"

		if getStatus.Decision == "REJ" {
			decision = "REJECTED"
		}

		if getStatus.Decision == "CAN" {
			decision = "CANCELED"
		}

		switch getStatus.SourceDecision {
		case "ACT":
			step = "ACTIVATION"
			digisign, _ := u.repository.GetLinkTrxDegisign(prospectID, "ACTIVATION")
			if digisign.Link != "" {
				activationUrl = digisign.Link
			}

		case "RGS":
			step = "REGISTER"
		case "SID":
			step = "SIGN_DOC"
			digisign, _ := u.repository.GetLinkTrxDegisign(prospectID, "SIGN_DOC")
			if getStatus.Decision == "APR" {
				decision = "APPROVED"
			}
			if digisign.Link != "" {
				signDocUrl = digisign.Link
			}

		case "SND":
			step = "SEND_DOC"
		}

		data = response.DataDigisignCheck{
			ProspectID:    prospectID,
			Step:          step,
			Decision:      decision,
			ActivationUrl: activationUrl,
			SignDocUrl:    signDocUrl,
		}
	} else {
		err = fmt.Errorf(constant.DIGISIGN_RECORD_NOT_FOUND)
		return
	}

	return

}
