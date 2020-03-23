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

func NewDigisignRegistrationResponse() *digisignRegistrationResponse {
	return &digisignRegistrationResponse{}
}

func (r *digisignRegistrationResponse) Bind(byteResponse []byte) error {
	library := jsoniter.ConfigCompatibleWithStandardLibrary
	err := library.Unmarshal(byteResponse, r)
	if err != nil {
		return err
	}
	return nil
}
