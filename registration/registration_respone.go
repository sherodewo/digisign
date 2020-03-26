package registration

import jsoniter "github.com/json-iterator/go"

type digisignRegistrationResponse struct {
	JSONFile struct {
		Result          string `json:"result"`
		Notif           string `json:"notif"`
		Info            string `json:"info,omitempty"`
		RefTrx          string `json:"refTrx,omitempty"`
		KodeUser        string `json:"kode_user,omitempty"`
		EmailRegistered string `json:"email_registered,omitempty"`
		ExpiredAktivasi string `json:"expired_aktivasi,omitempty"`
	}
}

type mapDigisignResponse struct {
	ProspectID      string `json:"prospect_id"`
	Result          string `json:"result"`
	Notif           string `json:"notif"`
	Info            string `json:"info,omitempty"`
	RefTrx          string `json:"refTrx,omitempty"`
	KodeUser        string `json:"kode_user,omitempty"`
	EmailRegistered string `json:"email_registered,omitempty"`
	ExpiredAktivasi string `json:"expired_aktivasi,omitempty"`
}

func NewDigisignRegistrationResponse() *digisignRegistrationResponse {
	return &digisignRegistrationResponse{}
}

func (r *digisignRegistrationResponse) Bind(prospectID string, byteResponse []byte) (*mapDigisignResponse, error) {
	library := jsoniter.ConfigCompatibleWithStandardLibrary
	err := library.Unmarshal(byteResponse, r)
	if err != nil {
		return nil, err
	}
	mapResponse := mapDigisignResponse{
		ProspectID:      prospectID,
		Result:          r.JSONFile.Result,
		Notif:           r.JSONFile.Notif,
		Info:            r.JSONFile.Info,
		RefTrx:          r.JSONFile.RefTrx,
		KodeUser:        r.JSONFile.KodeUser,
		EmailRegistered: r.JSONFile.EmailRegistered,
		ExpiredAktivasi: r.JSONFile.ExpiredAktivasi,
	}
	return &mapResponse, nil
}
