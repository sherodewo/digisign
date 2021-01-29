package sign_document

import (
	"los-int-digisign/infrastructure/config/digisign"
	"os"

	"gopkg.in/resty.v1"
)

type losSignDocumentRequestCallbackRequest struct {
	ClientKey      string `json:"client_key"`
	Email          string `json:"email"`
	Result         string `json:"result"`
	DocumentID     string `json:"document_id"`
	StatusDocument string `json:"status_document"`
}

func NewLosSignDocumentCallbackRequest() *losSignDocumentRequestCallbackRequest {
	return &losSignDocumentRequestCallbackRequest{}
}

func (c *losSignDocumentRequestCallbackRequest) losSignDocumentRequestCallback(email string, result string,
	documentId string, statusDocument string) (res *resty.Response, err error) {

	c.ClientKey = os.Getenv("LOS_KEY")
	c.Email = email
	c.Result = result
	c.DocumentID = documentId
	c.StatusDocument = statusDocument

	client := resty.New()
	// client.SetDebug(true)
	resp, err := client.R().SetHeader("Content-Type", "application/json").
		SetBody(c).
		Post(os.Getenv("LOS_BASE_URL") + "/digisign/sign_doc")
	if err != nil {
		tags := map[string]string{
			"app.pkg":      "sign_document",
			"app.func":     "losSignDocumentRequestCallback",
			"app.action": "forwarding-callback",
			"app.process": "sign-doc",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DIGISIGN-API")
		return nil, err
	}

	return resp, nil
}
