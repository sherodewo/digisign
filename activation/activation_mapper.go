package activation

import "kpdigisign/model"

type Mapper struct {
	ID                 string `json:"id"`
	ActivationResultID string `json:"activation_result_id"`
	EmailUser          string `json:"email_user"`
	Result             string `json:"result"`
	Link               string `json:"link"`
}

func NewActivationMapper() *Mapper {
	return &Mapper{}
}
func (m *Mapper) Map(model model.Activation) *Mapper {
	m.ID = model.ID
	m.ActivationResultID = model.ActivationResultID
	m.EmailUser = model.EmailUser
	m.Result = model.ActivationResult.Result
	m.Link = model.ActivationResult.Link
	return m
}

func (m *Mapper) MapList(model []model.Activation) *[]Mapper {
	var length = len(model)
	serialized := make([]Mapper, length)

	for k, v := range model {
		serialized[k] = Mapper{
			ID:                 v.ID,
			ActivationResultID: v.ActivationResult.ID,
			EmailUser:          v.EmailUser,
			Result:             v.ActivationResult.Result,
			Link:               v.ActivationResult.Link,
		}
	}
	return &serialized
}
