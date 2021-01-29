package registration

import (
	"bytes"
	"encoding/json"

	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/utils"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/resty.v1"
)

type digisignRegistrationRequest struct {
	JSONFile struct {
		UserID              string `json:"userid"`
		Alamat              string `json:"alamat"`
		JenisKelamin        string `json:"jenis_kelamin"`
		Kecamatan           string `json:"kecamatan"`
		Kelurahan           string `json:"kelurahan"`
		KodePos             string `json:"kode-pos"`
		Kota                string `json:"kota"`
		Nama                string `json:"nama"`
		NoTelepon           string `json:"tlp"`
		TanggalLahir        string `json:"tgl_lahir"`
		Provinsi            string `json:"provinci"`
		Nik                 string `json:"idktp"`
		TempatLahir         string `json:"tmp_lahir"`
		Email               string `json:"email"`
		Npwp                string `json:"npwp"`
		Redirect            bool   `json:"redirect"`
		AsliRiRefVerifikasi string `json:"ref_verifikasi,omitempty"`
		DataVerifikasi      string `json:"data_verifikasi,omitempty"`
		ScoreSelfie         string `json:"score_selfie,omitempty"`
		Vnik                string `json:"vnik,omitempty"`
		Vnama               string `json:"vnama,omitempty"`
		VtanggalLahir       string `json:"vtgl_lahir,omitempty"`
		VtempatLahir        string `json:"vtmp_lahir,omitempty"`
	}
}

func NewDigisignRegistrationRequest() *digisignRegistrationRequest {
	return &digisignRegistrationRequest{}
}

func (dr *digisignRegistrationRequest) DigisignRegistration(userType string, byteKtp []byte, byteSelfie []byte,
	byteNpwp []byte, byteTtd []byte, dto Dto) (res *resty.Response, result string, notif string, reftrx string,
	jsonResponse string, kodeUser string, request string, err error) {

	//Mapping request
	dr.JSONFile.UserID = dto.UserID
	dr.JSONFile.Alamat = dto.Alamat
	dr.JSONFile.JenisKelamin = dto.JenisKelamin
	dr.JSONFile.Kecamatan = dto.Kecamatan
	dr.JSONFile.Kelurahan = dto.Kelurahan
	dr.JSONFile.KodePos = dto.KodePos
	dr.JSONFile.Kota = dto.Kota
	dr.JSONFile.Nama = dto.Nama
	dr.JSONFile.NoTelepon = dto.NoTelepon
	dr.JSONFile.TanggalLahir = dto.TanggalLahir
	dr.JSONFile.Provinsi = dto.Provinsi
	dr.JSONFile.Nik = dto.Nik
	dr.JSONFile.TempatLahir = dto.TempatLahir
	dr.JSONFile.Email = dto.Email
	dr.JSONFile.Npwp = dto.Npwp
	dr.JSONFile.Redirect = dto.Redirect

	if userType == "NEW" {
		//Data AsliRI
		dr.JSONFile.AsliRiRefVerifikasi = *dto.AsliRiRefVerifikasi
		var dataVerifikasi map[string]interface{}
		dataVerifikasi = map[string]interface{}{
			"name":       dto.AsliRiNama,
			"birthplace": dto.AsliRiTempatLahir,
			"birthdate":  dto.AsliRiTanggalLahir,
			"address":    dto.AsliRiAlamat,
		}
		jsonDataVerifikasi, _ := json.Marshal(dataVerifikasi)
		dr.JSONFile.DataVerifikasi = string(jsonDataVerifikasi)
		dr.JSONFile.Vnik = *dto.Vnik
		dr.JSONFile.Vnama = *dto.Vnama
		dr.JSONFile.VtanggalLahir = *dto.VtanggalLahir
		dr.JSONFile.VtempatLahir = *dto.VtempatLahir
	}
	drJson, err := json.Marshal(dr)

	client := resty.New()
	client.SetTimeout(time.Second * time.Duration(250))

	if byteTtd == nil && byteNpwp == nil {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer "+digisign.Token).
			SetFileReader("fotoktp", "ktp_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post(os.Getenv("DIGISIGN_BASE_URL") + "/REG-MITRA.html")
		if err != nil {
			tags := map[string]string{
				"app.pkg":       "registration",
				"app.func":      "DigisignRegistration",
				"app.action":  "call",
				"app.process":  "register",
				"app.condition": "Ttd_Npwp_Nil",
			}
			extra := map[string]interface{}{
				"message":     err.Error(),
				"prospect_id": dto.ProspectID,
			}

			digisign.SendToSentry(tags, extra, "DIGISIGN-API")
			return nil, "", "", "", "", "", string(drJson), nil

		}
		result = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
		notif = jsoniter.Get(resp.Body(), "JSONFile").Get("notif").ToString()
		reftrx = jsoniter.Get(resp.Body(), "JSONFile").Get("refTrx").ToString()
		return resp, result, notif, reftrx, resp.String(), kodeUser, string(drJson), err
	} else if byteNpwp == nil {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer "+digisign.Token).
			SetFileReader("fotoktp", "ktp_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFileReader("ttd", "ttd_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteTtd),
				bytes.NewReader(byteTtd)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post(os.Getenv("DIGISIGN_BASE_URL") + "/REG-MITRA.html")
		if err != nil {
			tags := map[string]string{
				"app.pkg":       "registration",
				"app.func":      "DigisignRegistration",
				"app.action":  "call",
				"app.process":  "register",
				"app.condition": "Npwp_Nil",
			}
			extra := map[string]interface{}{
				"message":     err.Error(),
				"prospect_id": dto.ProspectID,
			}

			digisign.SendToSentry(tags, extra, "DIGISIGN-API")
			return nil, "", "", "", "", "", string(drJson), nil
		}
		result = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
		notif = jsoniter.Get(resp.Body(), "JSONFile").Get("notif").ToString()
		reftrx = jsoniter.Get(resp.Body(), "JSONFile").Get("refTrx").ToString()
		return resp, result, notif, reftrx, resp.String(), kodeUser, string(drJson), err
	} else if byteTtd == nil {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer "+digisign.Token).
			SetFileReader("fotoktp", "ktp_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFileReader("fotonpwp", "npwp_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteNpwp),
				bytes.NewReader(byteNpwp)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post(os.Getenv("DIGISIGN_BASE_URL") + "/REG-MITRA.html")
		if err != nil {
			tags := map[string]string{
				"app.pkg":       "registration",
				"app.func":      "DigisignRegistration",
				"app.action":  "call",
				"app.process":  "register",
				"app.condition": "Ttd_Nil",
			}
			extra := map[string]interface{}{
				"message":     err.Error(),
				"prospect_id": dto.ProspectID,
			}

			digisign.SendToSentry(tags, extra, "DIGISIGN-API")
			return nil, "", "", "", "", "", string(drJson), nil
		}
		result = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
		notif = jsoniter.Get(resp.Body(), "JSONFile").Get("notif").ToString()
		reftrx = jsoniter.Get(resp.Body(), "JSONFile").Get("refTrx").ToString()
		return resp, result, notif, reftrx, resp.String(), kodeUser, string(drJson), err
	} else {
		resp, err := client.R().
			SetHeader("Content-Type", "multipart/form-data").
			SetHeader("Authorization", "Bearer "+digisign.Token).
			SetFileReader("fotoktp", "ktp_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteKtp),
				bytes.NewReader(byteKtp)).
			SetFileReader("fotodiri", "selfie_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteSelfie),
				bytes.NewReader(byteSelfie)).
			SetFileReader("fotonpwp", "npwp_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteNpwp),
				bytes.NewReader(byteNpwp)).
			SetFileReader("ttd", "ttd_"+dto.Nama+"."+utils.GetExtensionImageFromByte(byteTtd),
				bytes.NewReader(byteTtd)).
			SetFormData(map[string]string{
				"jsonfield": string(drJson),
			}).
			Post(os.Getenv("DIGISIGN_BASE_URL") + "/REG-MITRA.html")
		if err != nil {
			tags := map[string]string{
				"app.pkg":       "registration",
				"app.func":      "DigisignRegistration",
				"app.action":  "call",
				"app.process":  "register",
				"app.condition": "Full_Data",
			}
			extra := map[string]interface{}{
				"message":     err.Error(),
				"prospect_id": dto.ProspectID,
			}

			digisign.SendToSentry(tags, extra, "DIGISIGN-API")
			return nil, "", "", "", "", "", string(drJson), nil
		}
		result = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
		notif = jsoniter.Get(resp.Body(), "JSONFile").Get("notif").ToString()
		reftrx = jsoniter.Get(resp.Body(), "JSONFile").Get("refTrx").ToString()
		kodeUser = jsoniter.Get(resp.Body(), "JSONFile").Get("kode_user").ToString()
		return resp, result, notif, reftrx, resp.String(), kodeUser, string(drJson), err
	}
}
