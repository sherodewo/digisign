package send_document

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty"
	jsoniter "github.com/json-iterator/go"
	"time"
	"os"
)

type digisignSendDocRequest struct {
	JSONFile struct {
		UserID         string      `json:"userid"`
		DocumentID     string      `json:"document_id"`
		Payment        string      `json:"payment"`
		SendTo         interface{} `json:"send-to"`
		ReqSign        interface{} `json:"req-sign"`
		Redirect       bool        `json:"redirect"`
		SequenceOption bool        `json:"sequence_option"`
		SigningSeq     int         `json:"signing_seq"`
	}
}

func NewDigisignSendDocRequest() *digisignSendDocRequest {
	return &digisignSendDocRequest{}
}

func (dr *digisignSendDocRequest) DigisignSendDoc(byteFile []byte, dto Dto) (
	res *resty.Response, result string, notif string, reftrx string, jsonResponse string, err error) {
	dr.JSONFile.UserID = dto.UserID
	dr.JSONFile.DocumentID = dto.DocumentID
	dr.JSONFile.Payment = dto.Payment
	dr.JSONFile.Redirect = dto.Redirect
	dr.JSONFile.SequenceOption = dto.SequenceOption
	dr.JSONFile.SigningSeq = dto.SigningSeq
	dr.JSONFile.SendTo = dto.SendTo
	dr.JSONFile.ReqSign = dto.ReqSign

	drJson, err := json.Marshal(dr)

	client := resty.New()
	//client.SetDebug(true)
	client.SetTimeout(time.Second * time.Duration(250))
	resp, err := client.R().
		SetHeader("Content-Type", "multipart/form-data").
		SetHeader("Authorization", "Bearer "+os.Getenv("DIGISIGN_TOKEN")).
		SetFileReader("file", "file_"+dto.DocumentID+".pdf", bytes.NewReader(byteFile)).
		SetFormData(map[string]string{
			"jsonfield": string(drJson),
		}).
		Post(os.Getenv("DIGISIGN_BASE_URL") + "/SendDocMitraAT.html")

	if err != nil {
		return nil, "", "", "", "", nil
	}
	result = jsoniter.Get(resp.Body(), "JSONFile").Get("result").ToString()
	notif = jsoniter.Get(resp.Body(), "JSONFile").Get("notif").ToString()
	reftrx = jsoniter.Get(resp.Body(), "JSONFile").Get("refTrx").ToString()
	return resp, result, notif, reftrx, resp.String(), err
}
