package sign_document

import "kpdigisign/model"

type Mapper struct {
	ID                   string `json:"id"`
	SignDocumentResultID string `json:"sign_document_result_id"`
	EmailUser            string `json:"email_user"`
	Result               string `json:"result"`
	Link                 string `json:"link"`
}

func NewSignDocumentMapper() *Mapper {
	return &Mapper{}
}
func (m *Mapper) Map(model model.SignDocument) *Mapper {
	m.ID = model.ID
	m.SignDocumentResultID = model.SignDocumentResultID
	m.EmailUser = model.EmailUser
	m.Result = model.SignDocumentResult.Result
	m.Link = model.SignDocumentResult.Link
	return m
}

func (m *Mapper) MapList(model []model.SignDocument) *[]Mapper {
	var length = len(model)
	serialized := make([]Mapper, length)

	for k, v := range model {
		serialized[k] = Mapper{
			ID:                   v.ID,
			SignDocumentResultID: v.SignDocumentResultID,
			EmailUser:            v.EmailUser,
			Result:               v.SignDocumentResult.Result,
			Link:                 v.SignDocumentResult.Link,
		}
	}
	return &serialized
}
