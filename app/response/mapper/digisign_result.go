package mapper

import "kpdigisign/app/models"

type digisignResultMapper struct {
	ID     string `json:"id"`
	LosID  string `json:"los_id"`
	Result string `json:"result"`
	Notif  string `json:"notif"`
	RefTrx string `json:"ref_trx"`
}

func NewDigisignResultMapper() *digisignResultMapper {
	return &digisignResultMapper{}
}
func (us *digisignResultMapper) Map(digisignRegistrationResult models.DigisignResult) *digisignResultMapper {
	us.ID = digisignRegistrationResult.ID
	us.LosID = digisignRegistrationResult.LosID
	us.Result = digisignRegistrationResult.Result
	us.Notif = digisignRegistrationResult.Notif
	us.Notif = digisignRegistrationResult.RefTrx

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
			RefTrx: v.RefTrx,
		}
	}
	return serialized
}
