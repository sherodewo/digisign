package sign_document

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strconv"
)

type digisignSignDocumentRequest struct {
	JSONFile struct {
		UserID     string `json:"userid"`
		DocumentID string `json:"document_id"`
		EmailUser  string `json:"email_user"`
		ViewOnly   bool   `json:"view_only"`
	}
}

func NewDigisignSignDocumentRequest() *digisignSignDocumentRequest {
	return &digisignSignDocumentRequest{}
}

func (dr *digisignSignDocumentRequest) DigisignSignDocumentRequest(request Dto) (
	res *resty.Response, result string, link string, err error) {
	dr.JSONFile.UserID = request.UserID
	dr.JSONFile.EmailUser = request.EmailUser
	dr.JSONFile.DocumentID = request.DocumentID
	dr.JSONFile.ViewOnly = false
	drJson, err := json.Marshal(dr)
	if err != nil {
		return nil, "", "", err
	}
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	res, err = client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer "+os.Getenv("DIGISIGN_TOKEN")).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/gen/genSignPage.html")
	if err != nil {
		return nil, "", "", err
	}
	result = jsoniter.Get(res.Body(), "JSONFile").Get("result").ToString()
	link = jsoniter.Get(res.Body(), "JSONFile").Get("link").ToString()

	return res, result, link, err
}
