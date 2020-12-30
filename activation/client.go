package activation

import (
	"bytes"
	"encoding/json"

	"los-int-digisign/infrastructure/config/digisign"
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/resty.v1"
)

type digisignActivationRequest struct {
	JSONFile struct {
		UserID    string `json:"userid"`
		EmailUser string `json:"email_user"`
	}
}

func NewDigisignActivationRequest() *digisignActivationRequest {
	return &digisignActivationRequest{}
}

func (dr *digisignActivationRequest) ActivationDigisign(request Dto) (
	result *resty.Response, resultActivation string, link string, jsonRequest string, err error) {
	dr.JSONFile.UserID = request.UserID
	dr.JSONFile.EmailUser = request.EmailUser
	drJson, err := json.Marshal(dr)
	if err != nil {
		return nil, "", "", string(drJson), err
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
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/gen/genACTPage.html")
	if err != nil {
		return nil, "", "", string(drJson), err
	}
	resultActivation = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
	link = jsoniter.Get(resp.Body(), "JSONFile").Get("link").ToString()

	return resp, resultActivation, link, string(drJson), err
}
