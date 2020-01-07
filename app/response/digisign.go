package response

import jsoniter "github.com/json-iterator/go"

type disisignResponse struct {
	JsonFile jsonFile `json:"JSONFile"`
}

type jsonFile struct {
	Result          string `json:"result"`
	Notif           string `json:"notif"`
	EmailRegistered string `json:"email_registered,omitempty"`
}

//type digisignVerificationResponse struct {
//	JsonFile jsonFile `json:"JSONFile"`
//}

//type disisignErrorResponse struct {
//	Result string `json:"result"`
//	Notif  string `json:"notif"`
//}

//type jsonFile struct {
//	Result string `json:"result"`
//	Info   string `json:"info"`
//	Data   *data  `json:"data,omitempty"`
//}

//type data struct {
//	Name        bool   `json:"name"`
//	Birthplace  bool   `json:"birthplace"`
//	Birthdate   bool   `json:"birthdate"`
//	Address     string `json:"address"`
//	SelfieMatch bool   `json:"selfie_match"`
//}

func NewDigisignResponse() *disisignResponse {
	return &disisignResponse{}
}

//func NewDigisignErrorResponse() *disisignResponse {
//	return &disisignResponse{}
//}

func (c *disisignResponse) Bind(d []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(d, c); err != nil {
		return err
	}
	return nil
}
