package activation

import jsoniter "github.com/json-iterator/go"

type digisignActivationResponse struct {
	JSONFile struct {
		Result string `json:"result"`
		Notif  string `json:"notif,omitempty"`
		Info   string `json:"info,omitempty"`
		RefTrx string `json:"refTrx,omitempty"`
		Link   string `json:"link,omitempty"`
	}
}

type mapDigisignResponse struct {
	ProspectID string `json:"prospect_id"`
	Result     string `json:"result"`
	Notif      string `json:"notif,omitempty"`
	Info       string `json:"info,omitempty"`
	RefTrx     string `json:"refTrx,omitempty"`
	Link       string `json:"link,omitempty"`
}

func NewDigisignActivationResponse() *digisignActivationResponse {
	return &digisignActivationResponse{}
}

func (r *digisignActivationResponse) Bind(prospectId string, byteResponse []byte) (*mapDigisignResponse, error) {
	library := jsoniter.ConfigCompatibleWithStandardLibrary
	err := library.Unmarshal(byteResponse, r)
	if err != nil {
		return nil, err
	}
	mapResponse := mapDigisignResponse{
		ProspectID: prospectId,
		Result:     r.JSONFile.Result,
		Notif:      r.JSONFile.Notif,
		Info:       r.JSONFile.Info,
		RefTrx:     r.JSONFile.RefTrx,
		Link:       r.JSONFile.Link,
	}
	return &mapResponse, nil
}
