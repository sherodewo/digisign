package activation

import (
	"github.com/go-resty/resty"
	jsoniter "github.com/json-iterator/go"
	"os"
)

type losActivationRequestCallbackRequest struct {
	ClientKey string `json:"client_key"`
	Email     string `json:"email"`
	Result    string `json:"result"`
	Notif     string `json:"notif"`
}

func NewLosActivationCallbackRequest() *losActivationRequestCallbackRequest {
	return &losActivationRequestCallbackRequest{}
}

func (c *losActivationRequestCallbackRequest) losActivationRequestCallback(email string, result string,
	notif string) (res *resty.Response, code string, message string, err error) {

	c.ClientKey = os.Getenv("LOS_KEY")
	c.Email = email
	c.Result = result
	c.Notif = notif

	client := resty.New()
	client.SetDebug(true)
	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetBody(c).
		Post(os.Getenv("LOS_BASE_URL") + "/digisign/activation")
	if err != nil {
		return nil, "", "", err
	}
	code = jsoniter.Get(resp.Body(), "code").ToString()
	message = jsoniter.Get(resp.Body(), "message").ToString()

	return resp, code, message, nil

}
