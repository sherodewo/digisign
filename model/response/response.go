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
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Error     interface{} `json:"errors"`
	RequestID string      `json:"request_id"`
	Timestamp string      `json:"timestamp"`
}

type ImageDecodeResponse struct {
	Messages string `json:"messages"`
	Data     struct {
		Encode string `json:"encode"`
	} `json:"data"`
	Errors interface{} `json:"errors"`
	Code   string      `json:"code"`
}
