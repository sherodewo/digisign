package sign_document

import (
	"bytes"
	"encoding/json"

	"los-int-digisign/infrastructure/config/digisign"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/resty.v1"
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
	res *resty.Response, result string, link string, notif string, jsonRequest string, err error) {
	dr.JSONFile.UserID = request.UserID
	dr.JSONFile.EmailUser = request.EmailUser
	dr.JSONFile.DocumentID = request.DocumentID
	dr.JSONFile.ViewOnly = false
	drJson, err := json.Marshal(dr)
	if err != nil {
		return nil, "", "", "", string(drJson), err
	}
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	res, err = client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer "+digisign.Token).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/gen/genSignPage.html")
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "sign_document",
			"app.func": "DigisignSignDocumentRequest",
			"app.action":  "call",
			"app.process":  "sign-doc",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
			"doc_id": request.DocumentID,
		}

		digisign.SendToSentry(tags, extra, "DIGISIGN-API")
		return nil, "", "", "", string(drJson), err
	}
	result = jsoniter.Get(res.Body(), "JSONFile").Get("result").ToString()
	link = jsoniter.Get(res.Body(), "JSONFile").Get("link").ToString()
	notif = jsoniter.Get(res.Body(), "JSONFile").Get("notif").ToString()

	return res, result, link, notif, string(drJson), err
}
