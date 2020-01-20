package client

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/helpers"
	"kpdigisign/app/request"
	"strconv"
)

type digisignRegistrationRequest struct {
	request.DigisignRegistrationRequest
}

type digisignSendDocRequest struct {
	request.SendDocumentRequest
}

type downloadRequest struct {
	request.DownloadRequest
}

type activationRequest struct {
	request.ActivationRequest
}

type signDocumentRequest struct {
	request.SignDocumentRequest
}

func NewDigisignRegistrationRequest() *digisignRegistrationRequest {
	return &digisignRegistrationRequest{}
}
func NewDigisignSendDocRequest() *digisignSendDocRequest {
	return &digisignSendDocRequest{}
}
func NewDownloadRequest() *downloadRequest {
	return &downloadRequest{}
}
func NewActivationRequest() *activationRequest {
	return &activationRequest{}
}

func NewSignDocRequest() *signDocumentRequest {
	return &signDocumentRequest{}
}

func (dr *digisignRegistrationRequest) DigisignRegistration(userType string, byteKtp []byte, byteSelfie []byte,
	byteNpwp []byte, byteTtd []byte, losRequest request.LosRequest) (result *resty.Response, err error) {

	//Mapping request
	dr.JsonFile.UserID = "adminkreditplus@tandatanganku.com"
	dr.JsonFile.Alamat = losRequest.Alamat
	dr.JsonFile.JenisKelamin = losRequest.JenisKelamin
	dr.JsonFile.Kecamatan = losRequest.Kecamatan
	dr.JsonFile.Kelurahan = losRequest.Kelurahan
	dr.JsonFile.KodePos = losRequest.KodePos
	dr.JsonFile.Kota = losRequest.Kota
	dr.JsonFile.Nama = losRequest.Nama
	dr.JsonFile.NoTelepon = losRequest.NoTelepon
	dr.JsonFile.TanggalLahir = losRequest.TanggalLahir
	dr.JsonFile.Provinsi = losRequest.Provinsi
	dr.JsonFile.Nik = losRequest.Nik
	dr.JsonFile.TempatLahir = losRequest.TempatLahir
	dr.JsonFile.Email = losRequest.Email
	dr.JsonFile.Npwp = losRequest.Npwp
	dr.JsonFile.RegNumber = losRequest.RegNumber
	dr.JsonFile.Redirect = true

	if userType == "NEW" {
		//Data AsliRI
		dr.JsonFile.AsliRiRefVerifikasi = *losRequest.AsliRiRefVerifikasi
		var dataVerifikasi map[string]interface{}
		dataVerifikasi = map[string]interface{}{
			"name":       losRequest.AsliRiNama,
			"birthplace": losRequest.AsliRiTempatLahir,
			"birthdate":  losRequest.AsliRiTanggalLahir,
			"address":    losRequest.AsliRiAlamat,
		}
		jsonDataVerifikasi, _ := json.Marshal(dataVerifikasi)
		dr.JsonFile.DataVerifikasi = string(jsonDataVerifikasi)

		dr.JsonFile.Vnik = *losRequest.Vnik
		dr.JsonFile.Vnama = *losRequest.Vnama
		dr.JsonFile.VtanggalLahir = *losRequest.VtanggalLahir
		dr.JsonFile.VtempatLahir = *losRequest.VtempatLahir
	}

	drJson, err := json.Marshal(dr)

	client := resty.New()
	if byteTtd == nil && byteNpwp == nil {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
			SetFileReader("fotoktp", "ktp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post("https://api.tandatanganku.com/REG-MITRA.html")
		log.Info("Response :", resp.Body())

		return resp, err
	} else if byteNpwp == nil {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
			SetFileReader("fotoktp", "ktp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFileReader("ttd", "ttd_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteTtd),
				bytes.NewReader(byteTtd)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post("https://api.tandatanganku.com/REG-MITRA.html")
		log.Info("Response :", resp.Body())

		return resp, err
	} else if byteTtd == nil {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
			SetFileReader("fotoktp", "ktp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFileReader("fotonpwp", "npwp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteNpwp),
				bytes.NewReader(byteNpwp)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post("https://api.tandatanganku.com/REG-MITRA.html")
		log.Info("Response :", resp.Body())

		return resp, err
	} else {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
			SetFileReader("fotoktp", "ktp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFileReader("fotonpwp", "npwp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteNpwp),
				bytes.NewReader(byteNpwp)).
			SetFileReader("ttd", "ttd_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteTtd),
				bytes.NewReader(byteTtd)).

			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post("https://api.tandatanganku.com/REG-MITRA.html")
		log.Info("Response :", resp.Body())

		return resp, err
	}
}

func (dr *digisignSendDocRequest) DigisignSendDoc(byteFile []byte, losRequest request.LosSendDocumentRequest) (
	result *resty.Response, err error) {
	dr.JsonFile.UserID = "adminkreditplus@tandatanganku.com"
	dr.JsonFile.DocumentID = losRequest.DocumentID
	dr.JsonFile.Payment = losRequest.Payment
	dr.JsonFile.Redirect = true
	dr.JsonFile.SequenceOption = false

	sendTo := jsoniter.Get([]byte(losRequest.SendTo), "sendTo").GetInterface()
	reqSign := jsoniter.Get([]byte(losRequest.ReqSign), "reqSign").GetInterface()

	dr.JsonFile.SendTo = sendTo
	dr.JsonFile.ReqSign = reqSign
	drJson, err := json.Marshal(dr)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFileReader("file", "file_"+losRequest.DocumentID+".pdf", bytes.NewReader(byteFile)).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		Post("https://api.tandatanganku.com/SendDocMitraAT.html")
	log.Info("Response :", resp.String())

	return resp, err
}

func (dr *downloadRequest) Download(downloadRequest request.LosDownloadDocumentRequest) (result *resty.Response, file string, err error) {
	dr.JsonFile.UserID = "adminkreditplus@tandatanganku.com"
	dr.JsonFile.DocumentID = downloadRequest.DocumentID
	drJson, err := json.Marshal(dr)
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post("https://api.tandatanganku.com/DWMITRA64.html")
	log.Info("Response :", resp.String())
	base64File := jsoniter.Get(resp.Body(), "JSONFile").Get("file").ToString()

	return resp, base64File, err
}

func (dr *downloadRequest) DownloadFile(downloadRequest request.LosDownloadDocumentRequest) (result *resty.Response, err error) {
	dr.JsonFile.UserID = "adminkreditplus@tandatanganku.com"
	dr.JsonFile.DocumentID = downloadRequest.DocumentID
	drJson, err := json.Marshal(dr)
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post("https://api.tandatanganku.com/DWMITRA.html")
	return resp, err
}

func (dr *activationRequest) ActivationDigisign(request request.LosActivationRequest) (
	result *resty.Response, resultActivation string, link string, err error) {
	dr.JSONFile.UserID = "adminkreditplus@tandatanganku.com"
	dr.JSONFile.EmailUser = request.EmailUser
	drJson, err := json.Marshal(dr)
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post("https://api.tandatanganku.com/gen/genACTPage.html")

	log.Info("Response :", resp.String())
	resultActivation = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
	link = jsoniter.Get(resp.Body(), "JSONFile").Get("link").ToString()

	return resp, resultActivation, link, err
}

func (dr *signDocumentRequest) DigisignSignDocumentRequest(request request.LosSignDocumentRequest) (
	result *resty.Response, resultActivation string, link string, err error) {
	dr.JSONFile.UserID = "adminkreditplus@tandatanganku.com"
	dr.JSONFile.EmailUser = request.EmailUser
	dr.JSONFile.DocumentID = request.DocumentID
	dr.JSONFile.ViewOnly = false
	drJson, err := json.Marshal(dr)
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post("https://api.tandatanganku.com/gen/genSignPage.html")

	log.Info("Response :", resp.String())
	resultActivation = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
	link = jsoniter.Get(resp.Body(), "JSONFile").Get("link").ToString()

	return resp, resultActivation, link, err
}