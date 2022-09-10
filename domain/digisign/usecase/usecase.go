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
	"los-int-digisign/shared/httpclient"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	packages struct {
		usecase interfaces.Usecase
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

func NewPackages(usecase interfaces.Usecase) interfaces.Packages {
	return &packages{usecase: usecase}
}

func NewUsecase(repository interfaces.Repository, httpclient httpclient.HttpClient) interfaces.Usecase {
	return &usecase{
		repository: repository,
		httpclient: httpclient,
	}
}

func NewMultiUsecase(repository interfaces.Repository, httpclient httpclient.HttpClient) (interfaces.MultiUsecase, interfaces.Usecase) {

	usecase := NewUsecase(repository, httpclient)
	packages := NewPackages(usecase)

	return &multiUsecase{
		repository: repository,
		httpclient: httpclient,
		usecase:    usecase,
		packages:   packages,
	}, usecase
}

func (u multiUsecase) Register(req request.RegisterRequest) (err error) {

	return
}

func (u packages) GetRegisterPhoto(ktpUrl, selfieUrl, signatureUrl, npwpUrl, prospectID string) (ktpBase64, selfieBase64, signatureBase64, npwpBase64 string, err error) {

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

	photoNpwp := GetIsMedia(signatureUrl)

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

	return
}

func (u usecase) Activation(req request.ActivationRequest) (err error) {

	return
}

func (u usecase) SendDoc(req request.DownloadRequest) (err error) {

	return
}

func (u usecase) SignUseCase(req request.SignDocDto) (uploadRes response.MediaServiceResponse,err error) {
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

		// 1. Sign Document to Digisign
		signRes, err := u.SignDoc(req.ProspectID, data)
		if err != nil {
			return uploadRes, err
		}
		fmt.Println(signRes)
		// 2. Download Document to local
		downloadDto := request.DownloadRequest{
			JSONFile: request.DownloadDto{
				UserID:     req.UserID,
				DocumentID: req.DocumentID,
			},
		}
		fileName, err = u.DownloadDoc(req.ProspectID, downloadDto)
		if err != nil {
			return uploadRes, err
		}
	}else {
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

func (u usecase) SignDoc(prospectID string,req request.SignDocRequest) (resp response.SignDocResponse, err error) {
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

func (u usecase) DownloadDoc(prospectID string,req request.DownloadRequest) (name string, err error) {
	url := os.Getenv("DIGISIGN_BASE_URL") + os.Getenv("SIGN_DOCUMENT_URL")
	// Type belum ada di platform (Dummy DULU GAN)
	param := map[string]interface{}{
		"jsonfield": req,
	}
	header := map[string]string{
		"Content-Type":  "multipart/form-data",
		"Authorization": os.Getenv("Bearer ") + os.Getenv("DIGISIGN_TOKEN"),
	}
	restyResp, err := u.httpclient.DigiAPI(url, http.MethodGet, param, "", header, 30, prospectID)
	var respDownload response.DownloadResponse
	if restyResp != nil && http.StatusOK == restyResp.StatusCode() {
		if err := json.Unmarshal(restyResp.Body(), &respDownload); err != nil {
			return name, err
		}
	}
	dec, err := base64.StdEncoding.DecodeString(respDownload.JsonFile.File)
	if err != nil {
		panic(err)
	}
	name = "document_signed_"+prospectID+"_"+time.Now().String()
	f, err := os.Create(name)
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
