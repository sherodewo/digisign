package download_document

import (
	"bytes"
	"encoding/json"

	"los-int-digisign/infrastructure/config/digisign"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/resty.v1"
)

type digisignDownloadRequest struct {
	JSONFile struct {
		UserID     string `json:"userid"`
		DocumentID string `json:"document_id"`
	}
}

func NewDigisignDownloadRequest() *digisignDownloadRequest {
	return &digisignDownloadRequest{}
}

func (dr *digisignDownloadRequest) DownloadFileBase64(downloadRequest Dto) (result *resty.Response, file string, err error) {
	dr.JSONFile.UserID = downloadRequest.UserID
	dr.JSONFile.DocumentID = downloadRequest.DocumentID
	drJson, err := json.Marshal(dr)
	if err != nil {
		return nil, "", err
	}
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer "+digisign.Token).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/DWMITRA64.html")

	if err != nil {
		return nil, "", err
	}
	base64File := jsoniter.Get(resp.Body(), "JSONFile").Get("file").ToString()

	return resp, base64File, err
}

func (dr *digisignDownloadRequest) DownloadFile(downloadRequest Dto) (result *resty.Response, err error) {
	dr.JSONFile.UserID = downloadRequest.UserID
	dr.JSONFile.DocumentID = downloadRequest.DocumentID
	drJson, err := json.Marshal(dr)
	if err != nil {
		return nil, err
	}
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer "+digisign.Token).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/DWMITRA.html")
	return resp, err
}
