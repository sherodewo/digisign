package common_models

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
