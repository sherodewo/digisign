package usecase

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"los-int-digisign/domain/digisign/interfaces"

	"los-int-digisign/model/entity"
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
	"los-int-digisign/shared/constant"
	"los-int-digisign/shared/httpclient"
	"los-int-digisign/shared/utils"
	"time"

	"net/http"
	"os"
	"strconv"
	"strings"

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

		fmt.Println(digiResp)

		data = RegisterMappingResponse(digiResp, req.ProspectID)

		return
	}

	resp, err := u.httpclient.RegisterAPI(os.Getenv("REGISTER_URL"), param, header, constant.METHOD_POST, 30, dataFile, req.ProspectID)

	fmt.Println(string(resp.Body()))

	if err != nil {
		return
	}

	json.Unmarshal(resp.Body(), &digiResp)

	data = RegisterMappingResponse(digiResp, req.ProspectID)

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

func (u usecase) Activation(req request.ActivationRequest) (data response.DataActivationResponse, err error) {

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

		data = ActivationMappingResponse(digiResp, req.ProspectID)

		return
	}

	resp, err := u.httpclient.ActivationAPI(os.Getenv("ACTIVATION_URL"), constant.METHOD_POST, param, header, 30, req.ProspectID)

	fmt.Println(string(resp.Body()))

	if err != nil {
		return
	}

	json.Unmarshal(resp.Body(), &digiResp)

	data = ActivationMappingResponse(digiResp, req.ProspectID)

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

	json.Unmarshal(byteDecrypt, &activationCallback)

	if activationCallback.Result == constant.CODE_00 && activationCallback.Notif == constant.ACTIVATION_CALLLBACK_SUCCESS {

		var (
			dataCustomer entity.CustomerPersonal
			sendDoc      response.DataSendDocResponse
			signDoc      response.DataSignDocResponse
		)
		// 1. get Order id by email and nik
		dataCustomer, err = u.repository.GetCustomerPersonalByEmailAndNik(activationCallback.Email, activationCallback.NIK)
		if err != nil {
			return
		}

		err = u.repository.UpdateStatusDigisignActivation(dataCustomer.ProspectID)
		if err != nil {
			return
		}

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

		var details []entity.TrxDetail

		var nextStep interface{}

		statusProcess := "ONP"
		activity := "PRCD"
		decision := "PAS"
		nextStep = "SID"

		if sendDoc.Decision == constant.DECISION_REJECT {
			statusProcess = "FIN"
			activity = "STOP"
			decision = "REJ"
			nextStep = nil
		}

		details = append(details, entity.TrxDetail{
			ProspectID: dataCustomer.ProspectID, StatusProcess: statusProcess, Activity: activity, Decision: decision,
			RuleCode: sendDoc.Code, SourceDecision: "SND", NextStep: nextStep, CreatedBy: "SYSTEM", Info: string(info),
		})

		// sign doc
		if sendDoc.Decision == constant.DECISION_PASS {
			signDoc, err = u.usecase.SignDocument(request.JsonFileSign{
				UserID:     os.Getenv("DIGISIGN_USER_ID"),
				DocumentID: sendDoc.DocumentID,
				Email:      activationCallback.Email,
				ViewOnly:   false,
			}, dataCustomer.ProspectID)

			if err != nil {
				return
			}

			signStatusProcess := "ONP"
			signActivity := "PRCD"
			signDecision := "PAS"

			if sendDoc.Decision == constant.DECISION_REJECT {
				signStatusProcess = "FIN"
				signActivity = "STOP"
				signDecision = "REJ"
			}

			data = signDoc

			details = append(details, entity.TrxDetail{
				ProspectID: dataCustomer.ProspectID, StatusProcess: signStatusProcess, Activity: signActivity, Decision: signDecision,
				RuleCode: sendDoc.Code, SourceDecision: "SID", NextStep: nextStep, CreatedBy: "SYSTEM", Info: sendDoc.DocumentID + ".pdf",
			})

			err = u.repository.SaveTrx(details)

			if err != nil {
				return
			}

			return

		}

		data = response.DataSignDocResponse{
			ProspectID:     sendDoc.ProspectID,
			Code:           sendDoc.Code,
			Decision:       sendDoc.Decision,
			DecisionReason: sendDoc.DecisionReason,
		}

		err = u.repository.SaveTrx(details)

		if err != nil {
			return
		}

		return
		// return send_doc gagal
	}

	// return callback activation gagal
	return
}

func (u usecase) GeneratePK(prospectID string) (document []byte, docID string, agreementNo string, err error) {

	param, _ := json.Marshal(map[string]interface{}{
		"ProspectID": prospectID,
	})

	resp, _ := u.httpclient.DocumentAPI(os.Getenv("GENERATE_PK_URL"), constant.METHOD_POST, param, map[string]string{}, 30, prospectID)

	if resp.StatusCode() != 200 {
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

	return
}

func (u packages) SendDoc(req request.SendDoc) (data response.DataSendDocResponse, err error) {

	// 1. find data by email and nik
	sendDoc, err := u.repository.GetSendDocData(req.ProspectID)
	if err != nil {
		return
	}

	// 2. generate pdf
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
				Page: "9", Llx: "33", Lly: "445", Urx: "546", Ury: "283", Visible: "1",
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

		data.DocumentID = documentID
		data.AgreementNo = agreementNo

		return
	}

	resp, err := u.httpclient.SendDocAPI(os.Getenv("SEND_DOC_URL"), constant.METHOD_POST, param, header, 30, dataFile, req.ProspectID)

	fmt.Println(string(resp.Body()))

	if err != nil {
		return
	}

	json.Unmarshal(resp.Body(), &digiResponse)

	data = SendDocMappingResponse(digiResponse, req.ProspectID)

	data.DocumentID = documentID
	data.AgreementNo = agreementNo

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

func (u usecase) SignDocument(req request.JsonFileSign, prospectID string) (data response.DataSignDocResponse, err error) {

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

		data = SignDocumentMappingResponse(digiResp, prospectID)

		return
	}

	signDoc, err := u.httpclient.SignDocAPI(os.Getenv("SIGN_DOC_URL"), constant.METHOD_POST, param, header, 30, prospectID)

	fmt.Println(string(signDoc.Body()))
	if err != nil {
		return
	}

	json.Unmarshal(signDoc.Body(), &digiResp)

	data = SignDocumentMappingResponse(digiResp, prospectID)

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

func (u multiUsecase) SignCallback(msg string) (upload response.MediaServiceResponse, err error) {

	decodeValue, _ := base64.StdEncoding.DecodeString(msg)

	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(os.Getenv("DIGISIGN_AES_KEY")))

	var signCallback response.SignCallback

	json.Unmarshal(byteDecrypt, &signCallback)

	if signCallback.StatusDocument == constant.SIGN_DOC_COMPLETE && signCallback.Result == constant.CODE_00 {

		var data entity.TrxDetail

		data, _ = u.repository.GetCustomerPersonalByEmail(signCallback.DocumentID)

		if err != nil {
			return
		}

		var download string

		download, err = u.usecase.DownloadDoc(data.ProspectID, request.DownloadRequest{
			DocumentID: signCallback.DocumentID,
			UserID:     os.Getenv("DIGISIGN_USER_ID"),
		})

		if err != nil {
			return
		}

		upload, err = u.usecase.UploadDoc(data.ProspectID, download)

		if err != nil {
			return
		}

	}

	return
}

// func (u usecase) SignUseCase(req request.SignDocDto) (uploadRes response.MediaServiceResponse, err error) {
// 	// Check Dummy Setting
// 	var fileName string
// 	dummy := os.Getenv("DUMMY")

// 	if dummy != "ON" {
// 		data := request.JsonFile{
// 			JsonFile: request.JsonFileSign{
// 				UserID:     req.UserID,
// 				DocumentID: req.DocumentID,
// 				Email:      req.Email,
// 				ViewOnly:   req.ViewOnly,
// 			},
// 		}

// 		// 1. Find Document by email and document_id

// 		// 2. Sign Document to Digisign
// 		signRes, err := u.SignDoc(req.ProspectID, data)
// 		if err != nil {
// 			return uploadRes, err
// 		}
// 		fmt.Println(signRes)
// 		// 2. Download Document to local
// 		downloadDto := request.DownloadRequest{
// 			UserID:     req.UserID,
// 			DocumentID: req.DocumentID,
// 		}

// 		fileName, err = u.DownloadDoc(req.ProspectID, downloadDto)
// 		if err != nil {
// 			return uploadRes, err
// 		}
// 	} else {
// 		fileName = "dummy_file.pdf"
// 	}

// 	// 3. Upload Document to Platform
// 	uploadRes, err = u.UploadDoc(req.ProspectID, fileName)
// 	if err != nil {
// 		return
// 	}

// 	// 4. Delete Document on Local
// 	//defer os.Remove(fileName)
// 	return
// }

// func (u usecase) SignDoc(prospectID string, req request.JsonFile) (resp response.SignDocResponse, err error) {
// 	url := os.Getenv("DIGISIGN_BASE_URL") + os.Getenv("SIGN_DOCUMENT_URL")
// 	// Type belum ada di platform (Dummy DULU GAN)
// 	param := map[string]interface{}{
// 		"jsonfield": req,
// 	}
// 	header := map[string]string{
// 		"Content-Type":  "multipart/form-data",
// 		"Authorization": os.Getenv("Bearer ") + os.Getenv("DIGISIGN_TOKEN"),
// 	}
// 	restyResp, err := u.httpclient.DigiAPI(url, http.MethodPost, param, "", header, 30, prospectID)
// 	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
// 		if err := json.Unmarshal(restyResp.Body(), &resp); err != nil {
// 			return resp, err
// 		}
// 	}
// 	return
// }

func (u usecase) DownloadDoc(prospectID string, req request.DownloadRequest) (pdfBase64 string, err error) {

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

	restyResp, err := u.httpclient.DigiAPI(os.Getenv("DOWNLOAD_URL"), constant.METHOD_POST, param, "", header, 30, prospectID)
	var respDownload response.DownloadResponse
	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
		if err := json.Unmarshal(restyResp.Body(), &respDownload); err != nil {
			return pdfBase64, err
		}
	}
	dec, err := base64.StdEncoding.DecodeString(respDownload.JsonFile.File)
	if err != nil {
		panic(err)
	}
	pdfBase64 = "document_signed_" + prospectID + "_" + time.Now().String()
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
	url := os.Getenv("MEDIA_BASE_URL") + os.Getenv("MEDIA_UPLOAD_URL")
	// Type belum ada di platform (Dummy DULU GAN)
	param := map[string]string{
		"type":         "perjanjian",
		"reference_no": prospectID,
	}
	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": os.Getenv("MEDIA_CLIENT_KEY"),
	}
	restyResp, err := u.httpclient.MediaClient(url, http.MethodPost, param, fileName, header, 30, prospectID)
	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
		if err := json.Unmarshal(restyResp.Body(), &uploadResp); err != nil {
			return uploadResp, err
		}
	}
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

	image, err := u.httpclient.MediaClient(url+os.Getenv("MEDIA_PATH"), "GET", nil, "", header, timeOut, customerID)

	if image.StatusCode() != 200 || err != nil {
		err = errors.New("error")
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
