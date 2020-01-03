package client

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty"
	"kpdigisign/app/helpers"
	"kpdigisign/app/request"
)

type digisignRegistrationRequest struct {
	request.DigisignRegistrationRequest
}

func NewDigisignRegistrationRequest() *digisignRegistrationRequest {
	return &digisignRegistrationRequest{}
}

func NewDigisignRegistrationRequestRoAo() *digisignRegistrationRequest {
	return &digisignRegistrationRequest{}
}

func (dr *digisignRegistrationRequest) DigisignRegistration(byteKtp []byte, byteSelfie []byte, byteNpwp []byte, byteTtd []byte,
	losRequest request.LosRequest) (result *resty.Response, err error) {

	//Mapping request
	dr.JsonFile.UserID = losRequest.UserID
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

	//Data AsliRI
	dr.JsonFile.AsliRiRegNumber = losRequest.AsliRiRegNumber
	dr.JsonFile.AsliRiRefVerifikasi = losRequest.AsliRiRefVerifikasi
	dr.JsonFile.AsliRiNama = losRequest.AsliRiNama
	dr.JsonFile.AsliRiTempatLahir = losRequest.AsliRiTempatLahir
	dr.JsonFile.AsliRiTanggalLahir = losRequest.AsliRiTanggalLahir
	dr.JsonFile.AsliRiAlamat = losRequest.AsliRiAlamat
	dr.JsonFile.AsliRiSelfieSimilarity = losRequest.AsliRiSelfieSimilarity
	dr.JsonFile.BranchID = losRequest.BranchID
	dr.JsonFile.EmailBm = losRequest.EmailBm

	drJson, err := json.Marshal(dr)

	client := resty.New()
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFileReader("fotoktp", "ktp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteKtp),
			bytes.NewReader(byteKtp)).
		SetFileReader("fotodiri", "selfie_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteSelfie),
			bytes.NewReader(byteSelfie)).
		SetFileReader("fotonpwp", "npwp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteNpwp),
			bytes.NewReader(byteSelfie)).
		SetFileReader("ttd", "ttd_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteNpwp),
			bytes.NewReader(byteSelfie)).

		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		Post("https://api.tandatanganku.com/REG-MITRA.html")

	//jsonFileResponse := jsoniter.Get(resp.Body(), "JSONFile")
	return resp, err
}

func (dr *digisignRegistrationRequest) DigisignRegistrationRoAo(byteKtp []byte, byteSelfie []byte, byteNpwp []byte, byteTtd []byte,
	losRequest request.LosRequest) (result *resty.Response, err error) {

	//Mapping request
	dr.JsonFile.UserID = losRequest.UserID
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

	drJson, err := json.Marshal(dr)

	client := resty.New()
	client.SetDebug(true)
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer WYm4d97LUaa7khMabTNJ9imwQEe87KDxRajcV8a3PvEonyAe14orOe4iGqpUYN").
		SetFileReader("fotoktp", "ktp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteKtp),
			bytes.NewReader(byteKtp)).
		SetFileReader("fotodiri", "selfie_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteSelfie),
			bytes.NewReader(byteSelfie)).
		SetFileReader("fotonpwp", "npwp_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteNpwp),
			bytes.NewReader(byteSelfie)).
		SetFileReader("ttd", "ttd_"+losRequest.Nama+"."+helpers.GetExtensionImageFromByte(byteNpwp),
			bytes.NewReader(byteSelfie)).

		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		Post("https://api.tandatanganku.com/REG-MITRA.html")

	//jsonFileResponse := jsoniter.Get(resp.Body(), "JSONFile")
	return resp, err
}