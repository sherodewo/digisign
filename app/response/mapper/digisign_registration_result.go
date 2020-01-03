package mapper

import "kpdigisign/app/models"

type digisignRegistrationResultMapper struct {
	ID       	string 		`json:"id"`
	LosRequestID string 	`json:"los_request_id"`
	Result  	string    	`json:"result"`
	Notif  		string    	`json:"notif"`
	Name  		bool    	`json:"name"`
	BirthPlace  bool    	`json:"birth_place"`
	BirthDate  	bool    	`json:"birth_date"`
	Address  	string    	`json:"address"`
	Info  		string    	`json:"info"`
}

func NewDigisignRegistrationResultMapper() *digisignRegistrationResultMapper {
	return &digisignRegistrationResultMapper{}
}
func (us *digisignRegistrationResultMapper) Map(digisignRegistrationResult models.DigisignRegistrationResult) *digisignRegistrationResultMapper {
	us.ID = digisignRegistrationResult.ID
	us.LosRequestID = digisignRegistrationResult.LosRequestID
	us.Result = digisignRegistrationResult.Result
	us.Notif = digisignRegistrationResult.Notif
	us.Name = digisignRegistrationResult.Name
	us.BirthPlace = digisignRegistrationResult.BirthPlace
	us.BirthDate = digisignRegistrationResult.BirthDate
	us.Address = digisignRegistrationResult.Address
	us.Info = digisignRegistrationResult.Info

	return us
}

func (us *digisignRegistrationResultMapper) MapList(digisignRegistrationResult []models.DigisignRegistrationResult) interface{} {
	var length = len(digisignRegistrationResult)
	serialized := make([]digisignRegistrationResultMapper, length, length)

	for k, v := range digisignRegistrationResult {
		serialized[k] = digisignRegistrationResultMapper{
			ID:       	v.ID,
			Result: 	v.Result,
			Notif:    	v.Notif,
			Name:    	v.Name,
			BirthPlace: v.BirthPlace,
			BirthDate:  v.BirthDate,
			Address:    v.Address,
			Info:    	v.Info,
		}
	}
	return serialized
}
