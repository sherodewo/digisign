package response

type RegisterResponse struct {
	JsonFile struct {
		Result          string `json:"result"`
		Notif           string `json:"notif"`
		RefTrx          string `json:"refTrx"`
		Info            string `json:"info"`
		EmailRegistered string `json:"email_registered"`
	} `json:"JSONFile"`
}

type ActivationResponse struct {
	JsonFile struct {
		Result string `json:"result"`
		Link   string `json:"link"`
	} `json:"JSONFile"`
}

type SendDocResponse struct {
	JsonFile struct {
		Result string `json:"result"`
		Notif  string `json:"notif"`
		RefTrx string `json:"refTrx"`
		Info   string `json:"info"`
	} `json:"JSONFile"`
}

type SignDocResponse struct {
	JsonFile struct {
		Result string `json:"result"`
		Notif  string `json:"notif"`
		Link   string `json:"link"`
	} `json:"JSONFile"`
}

type DownloadResponse struct {
	JsonFile struct {
		Result string `json:"result"`
		Notif  string `json:"notif"`
		File   string `json:"file"`
	} `json:"JSONFile"`
}

type MediaServiceResponse struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Errors  interface{} `json:"errors"`
	Data    struct {
		Type        string `json:"type"`
		ReferenceNo string `json:"reference_no"`
		MediaURL    string `json:"media_url"`
		Path        string `json:"path"`
	} `json:"data"`
}

type ImageDecodeResponse struct {
	Messages string `json:"messages"`
	Data     struct {
		Encode string `json:"encode"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
	Code   string      `json:"code"`
}

type Api struct {
	Message    string      `json:"messages"`
	Errors     interface{} `json:"errors"`
	Data       interface{} `json:"data"`
	ServerTime string      `json:"server_time"`
}

type ErrorValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ActivationCallbackResponse struct {
	Result string `json:"result"`
	Notif  string `json:"notif"`
	Email  string `json:"email_user"`
	NIK    string `json:"nik"`
}

type SignCallback struct {
	DocumentID     string `json:"document_id"`
	StatusDocument string `json:"status_document"`
	Result         string `json:"result"`
	Notif          string `json:"notif"`
	Email          string `json:"email_user"`
}

type SignResponse struct {
	ProspectID string `json:"prospect_id"`
	Url        string `json:"url"`
}

type DataRegisterResponse struct {
	ProspectID     string `json:"prospect_id"`
	Code           string `json:"code"`
	Decision       string `json:"decision"`
	DecisionReason string `json:"decision_reason"`
	IsRegistered   bool   `json:"is_registered"`
}

type DataActivationResponse struct {
	ProspectID     string `json:"prospect_id"`
	Code           string `json:"code"`
	Decision       string `json:"decision"`
	DecisionReason string `json:"decision_reason"`
	Link           string `json:"activation_url"`
}

type DataSendDocResponse struct {
	ProspectID     string `json:"prospect_id"`
	Code           string `json:"code"`
	Decision       string `json:"decision"`
	DecisionReason string `json:"decision_reason"`
	DocumentID     string `json:"document_id"`
	AgreementNo    string `json:"agreement_no"`
}

type DataSignDocResponse struct {
	ProspectID     string `json:"prospect_id"`
	Code           string `json:"code"`
	Decision       string `json:"decision"`
	DecisionReason string `json:"decision_reason"`
	Link           string `json:"sign_doc_url"`
}

type DocumentGenerateResponse struct {
	DocumentPK  string `json:"doc_pk"`
	DocumentID  string `json:"doc_id"`
	AgreementNo string `json:"agreement_no"`
}

type SendDocInfo struct {
	DocumentID  string `json:"document_id"`
	AgreementNo string `json:"AgreementNo"`
}

type DataDigisignCheck struct {
	ProspectID    string      `json:"prospect_id"`
	Step          string      `json:"step"`
	Decision      string      `json:"decision"`
	ActivationUrl interface{} `json:"activation_url"`
	SignDocUrl    interface{} `json:"sign_doc_url"`
}
