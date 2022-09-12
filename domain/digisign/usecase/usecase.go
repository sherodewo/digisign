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
	"net/http"
	"os"
	"strconv"
	"strings"
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

func (u usecase) Activation(req request.ActivationRequest) (res response.ActivationResponse, err error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return
	}

	params := map[string]string{
		"jsonfield": string(jsonData),
	}

	resp, err := u.httpclient.DigisignAPI(os.Getenv("DIGISIGN_BASE_URL")+"/gen/genACTPage.html", params, map[string]string{}, constant.METHOD_POST, 10, false, nil)
	if err != nil || resp.StatusCode() != http.StatusOK {
		err = errors.New("error while do activation: " + err.Error())
		return
	}

	body := resp.Body()
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = errors.New("error while unmarshal activation response: " + err.Error())
		return
	}

	return
}

func (u usecase) SendDoc(req request.DownloadRequest) (err error) {

	return
}

func (u multiUsecase) SignDoc(req request.SignDocRequest) (err error) {

	return
}

func (usecase) UploadDoc() (err error) {

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

	image, err := u.httpclient.MediaClient(url+os.Getenv("MEDIA_PATH"), "GET", nil, header, timeOut, customerID)

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
