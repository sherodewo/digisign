package sign_document

import jsoniter "github.com/json-iterator/go"

type digisignSignDocumentResponse struct {
	JSONFile struct {
		Result string `json:"result"`
		Notif  string `json:"notif,omitempty"`
		Info   string `json:"info,omitempty"`
		RefTrx string `json:"refTrx,omitempty"`
		Link   string `json:"refTrx,omitempty"`
	}
}

type mapDigisignResponse struct {
	Result string `json:"result"`
	Notif  string `json:"notif,omitempty"`
	Info   string `json:"info,omitempty"`
	RefTrx string `json:"refTrx,omitempty"`
	Link   string `json:"refTrx,omitempty"`
}

func NewDigisignSignDocumentResponse() *digisignSignDocumentResponse {
	return &digisignSignDocumentResponse{}
}

func (r *digisignSignDocumentResponse) Bind(byteResponse []byte) (*mapDigisignResponse, error) {
	library := jsoniter.ConfigCompatibleWithStandardLibrary
	err := library.Unmarshal(byteResponse, r)
	if err != nil {
		return nil, err
	}
	mapResponse := mapDigisignResponse{
		Result: r.JSONFile.Result,
		Notif:  r.JSONFile.Notif,
		Info:   r.JSONFile.Info,
		RefTrx: r.JSONFile.RefTrx,
		Link:   r.JSONFile.Link,
	}
	return &mapResponse, nil
}
