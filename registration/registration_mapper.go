package registration

import "kpdigisign/model"

type Mapper struct {
	ID                   string `json:"id"`
	RegistrationResultID string `json:"registration_result_id"`
	Result               string `json:"result"`
	Notif                string `json:"notif"`
	RefTrx               string `json:"ref_trx"`
}

func NewRegistrationMapper() *Mapper {
	return &Mapper{}
}
func (m *Mapper) Map(model model.Registration) *Mapper {
	m.ID = model.ID
	m.RegistrationResultID = model.RegistrationResultID
	m.Result = model.RegistrationResult.Result
	m.Notif = model.RegistrationResult.Notif
	m.RefTrx = model.RegistrationResult.RefTrx
	return m
}

func (m *Mapper) MapList(model []model.Registration) *[]Mapper {
	var length = len(model)
	serialized := make([]Mapper, length)

	for k, v := range model {
		serialized[k] = Mapper{
			ID:                   v.ID,
			RegistrationResultID: v.RegistrationResult.ID,
			Result:               v.RegistrationResult.Result,
			Notif:                v.RegistrationResult.Notif,
			RefTrx:               v.RegistrationResult.RefTrx,
		}
	}
	return &serialized
}
