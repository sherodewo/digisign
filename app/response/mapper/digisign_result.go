package mapper

import "kpdigisign/app/models"

type digisignResultMapper struct {
	ID     string `json:"id"`
	LosID  string `json:"los_id"`
	Result string `json:"result"`
	Notif  string `json:"notif"`
	//Name       bool   `json:"name"`
	//BirthPlace bool   `json:"birth_place"`
	//BirthDate  bool   `json:"birth_date"`
	//Address    string `json:"address"`
	//Info       string `json:"info"`
}

func NewDigisignResultMapper() *digisignResultMapper {
	return &digisignResultMapper{}
}
func (us *digisignResultMapper) Map(digisignRegistrationResult models.DigisignResult) *digisignResultMapper {
	//us.ID = digisignRegistrationResult.ID
	us.LosID = digisignRegistrationResult.LosID
	us.Result = digisignRegistrationResult.Result
	us.Notif = digisignRegistrationResult.Notif
	//us.Name = digisignRegistrationResult.Name
	//us.BirthPlace = digisignRegistrationResult.BirthPlace
	//us.BirthDate = digisignRegistrationResult.BirthDate
	//us.Address = digisignRegistrationResult.Address
	//us.Info = digisignRegistrationResult.Info

	return us
}

func (us *digisignResultMapper) MapList(digisignRegistrationResult []models.DigisignResult) interface{} {
	var length = len(digisignRegistrationResult)
	serialized := make([]digisignResultMapper, length, length)

	for k, v := range digisignRegistrationResult {
		serialized[k] = digisignResultMapper{
			ID:     v.ID,
			LosID:  v.LosID,
			Result: v.Result,
			Notif:  v.Notif,
			//Name:       v.Name,
			//BirthPlace: v.BirthPlace,
			//BirthDate:  v.BirthDate,
			//Address:    v.Address,
			//Info:       v.Info,
		}
	}
	return serialized
}
