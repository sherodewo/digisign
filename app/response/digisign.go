package response

import jsoniter "github.com/json-iterator/go"

type checkLosRequest struct {
	Result 			string 	`json:"result"`
	Info 			string 	`json:"info"`
	EmailRegistered string 	`json:"email_registered"`
	Name 			bool 	`json:"name"`
	Birthplace 		bool 	`json:"birthplace"`
	Birthdate 		bool 	`json:"birthdate"`
	Address 		string 	`json:"address"`
	SelfieMatch 	bool 	`json:"selfie_match"`
	Data 			string 	`json:"data"`
}

func NewLosRespone() *checkLosRequest  {
	return &checkLosRequest{}
}

func (c *checkLosRequest) Bind(d []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err := json.Unmarshal(d, c); err != nil {
		return err
	}
	return nil
}
