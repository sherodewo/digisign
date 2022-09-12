package usecase

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"los-int-digisign/domain/digisign/interfaces"
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

func NewMultiUsecase(repository interfaces.Repository, httpclient httpclient.HttpClient) (interfaces.MultiUsecase, interfaces.Usecase) {

	usecase := NewUsecase(repository, httpclient)
	packages := NewPackages(repository, usecase, httpclient)

	return &multiUsecase{
		repository: repository,
		httpclient: httpclient,
		usecase:    usecase,
		packages:   packages,
	}, usecase
}

func (u multiUsecase) Register(req request.Register) (data response.RegisterResponse, err error) {

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
		Provinci:   req.Provinci,
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

	resp, err := u.httpclient.RegisterAPI(os.Getenv("REGISTER_URL"), param, header, constant.METHOD_POST, 30, dataFile, req.ProspectID)

	if err != nil {
		return
	}

	json.Unmarshal(resp.Body(), &data)

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

func (u usecase) Activation(req request.ActivationRequest) (res response.ActivationResponse, err error) {
	url := os.Getenv("DIGISIGN_BASE_URL") + os.Getenv("DIGISIGN_ACTIVATION_URL")

	params := map[string]interface{}{
		"jsonfield": req,
	}

	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": os.Getenv("Bearer ") + os.Getenv("DIGISIGN_TOKEN"),
	}

	resp, err := u.httpclient.DigiAPI(url, http.MethodPost, params, "", header, 30, "")
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = errors.New("error while do activation: " + err.Error())
		return
	}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		err = errors.New("error while unmarshal activation response: " + err.Error())
		return
	}

	return
}

func (u multiUsecase) ActivationRedirect(msg string) (data interface{}, err error) {

	decodeValue, _ := base64.StdEncoding.DecodeString(msg)

	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(os.Getenv("DIGISIGN_AES_KEY")))

	var activationCallback response.ActivationCallbackResponse

	json.Unmarshal(byteDecrypt, &activationCallback)

	if activationCallback.Result == "00" && activationCallback.Notif == "Proses Aktivasi Berhasil" {

		var (
			prospectID string
			sendDoc    response.SendDocResponse
			signDoc    response.SignDocResponse
			documentId string
		)
		// 1. get Order id by email and nik

		// send doc
		sendDoc, documentId, err = u.packages.SendDoc(prospectID)

		if err != nil {
			return
		}

		// sign doc
		if sendDoc.JsonFile.Result == "00" {
			signDoc, err = u.usecase.SignDocV2(request.JsonFileSign{
				UserID:     os.Getenv("DIGISIGN_USER_ID"),
				DocumentID: documentId,
				Email:      activationCallback.Email,
				ViewOnly:   false,
			}, prospectID)

			if err != nil {
				return nil, err
			}

			data = signDoc.JsonFile

			return

		}

		// return send_doc gagal
	}

	// return callback activation gagal

	return
}

func (u usecase) GeneratePK(prospectID string) (document []byte, docID string, err error) {

	// 1. Find data by prospectid

	return
}

func (u packages) SendDoc(prospectID string) (data response.SendDocResponse, documentID string, err error) {

	// 1. find data by email and nik

	// 2. generate pdf
	document, documentID, err := u.usecase.GeneratePK(prospectID)
	if err != nil {
		return
	}

	dataFile := request.DataFile{
		DocumentPK: document,
		DocumentID: documentID,
	}

	jsonFile := request.SendDocRequest{
		UserID:         os.Getenv("DIGISIGN_USER_ID"),
		DocumentID:     strconv.FormatInt(time.Now().Unix(), 10),
		Payment:        os.Getenv("PAYMENT_METHOD"),
		Redirect:       true,
		Branch:         "", // from db
		SequenceOption: false,
		SendTo: []request.SendTo{
			{
				Name:  "", // from legal name customer
				Email: "", // from email customer
			},
			{
				Name:  "", // from data bm
				Email: "", // from data bm
			},
		},
		ReqSign: []request.ReqSign{
			{
				Name:    "", // from data bm
				Email:   "", // from data bm
				User:    "prf1",
				Page:    "1",
				Llx:     "323",
				Lly:     "135",
				Urx:     "420",
				Ury:     "184",
				Visible: "1",
			},
			{
				Name:    "", // from data customer
				Email:   "", // email customer
				User:    "ttd1",
				Page:    "1",
				Llx:     "458",
				Lly:     "135",
				Urx:     "557",
				Ury:     "184",
				Visible: "1",
			},
			{
				Name:    "", // from data customer
				Email:   "", // email customer
				User:    "ttd2",
				Page:    "5",
				Llx:     "70",
				Lly:     "356.7",
				Urx:     "183",
				Ury:     "457.5",
				Visible: "1",
			},
			{
				Name:    "", // from data bm
				Email:   "", // from data bm
				User:    "prf2",
				Page:    "5",
				Llx:     "428.4",
				Lly:     "356.7",
				Urx:     "541.4",
				Ury:     "457.5",
				Visible: "1",
			},
			{
				Name:    "", // from data bm
				Email:   "", // from data bm
				User:    "prf3",
				Page:    "7",
				Llx:     "33",
				Lly:     "448.6",
				Urx:     "126.7",
				Ury:     "495.4",
				Visible: "1",
			},
			{
				Name:    "", // from data customer
				Email:   "", // email customer
				User:    "ttd3",
				Page:    "7",
				Llx:     "457",
				Lly:     "448.6",
				Urx:     "580",
				Ury:     "495.4",
				Visible: "1",
			},
			{
				Name:    "", // from data customer
				Email:   "", // email customer
				User:    "ttd4",
				Page:    "9",
				Llx:     "71.3",
				Lly:     "251",
				Urx:     "130",
				Ury:     "283",
				Visible: "1",
			},
			{
				Name:    "", // from data bm
				Email:   "", // from data bm
				User:    "prf4",
				Page:    "9",
				Llx:     "33",
				Lly:     "445",
				Urx:     "546",
				Ury:     "283",
				Visible: "1",
			},
			{
				Name:    "", // from data customer
				Email:   "", // email customer
				User:    "ttd5",
				Page:    "10",
				Llx:     "31",
				Lly:     "180",
				Urx:     "118",
				Ury:     "276,5",
				Visible: "1",
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

	resp, err := u.httpclient.SendDocAPI(os.Getenv("SEND_DOC_URL"), constant.METHOD_POST, param, header, 30, dataFile, prospectID)

	if err != nil {
		return
	}

	json.Unmarshal(resp.Body(), &data)

	return
}

func (u usecase) SignDocV2(req request.JsonFileSign, prospectID string) (data response.SignDocResponse, err error) {

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

	signDoc, err := u.httpclient.SignDocAPI(os.Getenv("SIGN_DOC_URL"), constant.METHOD_POST, param, header, 30, prospectID)

	if err != nil {
		return
	}

	json.Unmarshal(signDoc.Body(), &data)

	return
}

func (u multiUsecase) SignCallback(msg string) (upload response.MediaServiceResponse, err error) {

	decodeValue, _ := base64.StdEncoding.DecodeString(msg)

	byteDecrypt := utils.AesDecrypt(decodeValue, []byte(os.Getenv("DIGISIGN_AES_KEY")))

	var signCallback response.SignCallback

	json.Unmarshal(byteDecrypt, &signCallback)

	if signCallback.StatusDocument == "complete" && signCallback.Result == "00" {

		// 1. Find data by document id dan email
		var (
			prospectID string
			download   string
		)

		download, err = u.usecase.DownloadDoc(prospectID, request.DownloadRequest{
			DocumentID: signCallback.DocumentID,
			UserID:     os.Getenv("DIGISIGN_USER_ID"),
		})

		if err != nil {
			return
		}

		upload, err = u.usecase.UploadDoc(prospectID, download)

		if err != nil {
			return
		}

	}

	return
}

func (u usecase) SignUseCase(req request.SignDocDto) (uploadRes response.MediaServiceResponse, err error) {
	// Check Dummy Setting
	var fileName string
	dummy := os.Getenv("DUMMY")

	if dummy != "ON" {
		data := request.SignDocRequest{
			JsonFile: request.JsonFileSign{
				UserID:     req.UserID,
				DocumentID: req.DocumentID,
				Email:      req.Email,
				ViewOnly:   req.ViewOnly,
			},
		}

		// 1. Find Document by email and document_id

		// 2. Sign Document to Digisign
		signRes, err := u.SignDoc(req.ProspectID, data)
		if err != nil {
			return uploadRes, err
		}
		fmt.Println(signRes)
		// 2. Download Document to local
		downloadDto := request.DownloadRequest{
			UserID:     req.UserID,
			DocumentID: req.DocumentID,
		}

		fileName, err = u.DownloadDoc(req.ProspectID, downloadDto)
		if err != nil {
			return uploadRes, err
		}
	} else {
		fileName = "dummy_file.pdf"
	}

	// 3. Upload Document to Platform
	uploadRes, err = u.UploadDoc(req.ProspectID, fileName)
	if err != nil {
		return
	}

	// 4. Delete Document on Local
	//defer os.Remove(fileName)
	return
}

func (u usecase) SignDoc(prospectID string, req request.SignDocRequest) (resp response.SignDocResponse, err error) {
	url := os.Getenv("DIGISIGN_BASE_URL") + os.Getenv("SIGN_DOCUMENT_URL")
	// Type belum ada di platform (Dummy DULU GAN)
	param := map[string]interface{}{
		"jsonfield": req,
	}
	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": os.Getenv("Bearer ") + os.Getenv("DIGISIGN_TOKEN"),
	}
	restyResp, err := u.httpclient.DigiAPI(url, http.MethodPost, param, "", header, 30, prospectID)
	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
		if err := json.Unmarshal(restyResp.Body(), &resp); err != nil {
			return resp, err
		}
	}
	return
}

func (u usecase) DownloadDoc(prospectID string, req request.DownloadRequest) (pdfBase64 string, err error) {

	param := map[string]interface{}{
		"jsonfield": req,
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
		"type":         "ePO",
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
