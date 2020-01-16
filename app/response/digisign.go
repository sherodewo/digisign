package response

import jsoniter "github.com/json-iterator/go"

type digisignResponse struct {
	JsonFile jsonFile `json:"JSONFile"`
}

type jsonFile struct {
	Result string `json:"result"`
	Notif  string `json:"notif"`
	RefTrx string `json:"refTrx"`
}

func NewDigisignResponse() *digisignResponse {
	return &digisignResponse{}
}

func (c *digisignResponse) Bind(d []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(d, c); err != nil {
		return err
	}
	return nil
}
