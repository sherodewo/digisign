package send_document

import "kpdigisign/model"

type Mapper struct {
	ID                   string `json:"id"`
	SendDocumentResultID string `json:"send_document_result_id"`
	Result               string `json:"result"`
	Notif                string `json:"notif"`
	RefTrx               string `json:"ref_trx"`
}

func NewSendDocumentMapper() *Mapper {
	return &Mapper{}
}
func (m *Mapper) Map(model model.SendDocument) *Mapper {
	m.ID = model.ID
	m.SendDocumentResultID = model.SendDocumentResultID
	m.Result = model.SendDocumentResult.Result
	m.Notif = model.SendDocumentResult.Notif
	m.RefTrx = model.SendDocumentResult.RefTrx
	return m
}

func (m *Mapper) MapList(model []model.SendDocument) *[]Mapper {
	var length = len(model)
	serialized := make([]Mapper, length)

	for k, v := range model {
		serialized[k] = Mapper{
			ID:                   v.ID,
			SendDocumentResultID: v.SendDocumentResult.ID,
			Result:               v.SendDocumentResult.Result,
			Notif:                v.SendDocumentResult.Notif,
			RefTrx:               v.SendDocumentResult.RefTrx,
		}
	}
	return &serialized
}
