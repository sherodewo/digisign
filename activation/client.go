package activation

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty"
	jsoniter "github.com/json-iterator/go"
	"os"
	"strconv"
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
	result *resty.Response, resultActivation string, link string, err error) {
	dr.JSONFile.UserID = os.Getenv("DIGISIGN_USER_ID")
	dr.JSONFile.EmailUser = request.EmailUser
	drJson, err := json.Marshal(dr)
	if err != nil {
		return nil, "", "", err
	}
	bs := []byte(strconv.Itoa(0))
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer "+os.Getenv("DIGISIGN_TOKEN")).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		SetFileReader("file", "file_", bytes.NewReader(bs)).
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/gen/genACTPage.html")
	if err != nil {
		return nil, "", "", err
	}
	resultActivation = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
	link = jsoniter.Get(resp.Body(), "JSONFile").Get("link").ToString()

	return resp, resultActivation, link, err
}
