package response

type RegisterResponse struct {
	JsonFile struct {
		Result string `json:"result"`
		Notif  string `json:"notif"`
		RefTrx string `json:"refTrx"`
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
	} `json:"JSONFile"`
}

type SignDocResponse struct {
	JsonFile struct {
		Result string `json:"result"`
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
type SignResponse struct {
	ProspectID string `json:"prospect_id"`
	Url        string `json:"url"`
}
