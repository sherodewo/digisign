package activation

import "los-int-digisign/model"

type Mapper struct {
	ID                 string `json:"id"`
	ActivationResultID string `json:"activation_result_id"`
	ProspectID         string `json:"prospect_id"`
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
	m.ProspectID = model.ProspectID
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
			ProspectID:         v.ProspectID,
			EmailUser:          v.EmailUser,
			Result:             v.ActivationResult.Result,
			Link:               v.ActivationResult.Link,
		}
	}
	return &serialized
}
