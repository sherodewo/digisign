package mapper

import "kpdigisign/app/models"

type documentResultMapper struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Result     string `json:"result"`
	Notif      string `json:"notif"`
}

func NewDocumentResultMapper() *documentResultMapper {
	return &documentResultMapper{}
}
func (us *documentResultMapper) Map(digisignDocumentResult models.DocumentResult) *documentResultMapper {
	us.ID = digisignDocumentResult.ID
	us.DocumentID = digisignDocumentResult.DocumentID
	us.Result = digisignDocumentResult.Result
	us.Notif = digisignDocumentResult.Notif
	return us
}

func (us *documentResultMapper) MapList(digisignDocumentResult []models.DocumentResult) interface{} {
	var length = len(digisignDocumentResult)
	serialized := make([]documentResultMapper, length, length)

	for k, v := range digisignDocumentResult {
		serialized[k] = documentResultMapper{
			ID:     v.ID,
			DocumentID:  v.DocumentID,
			Result: v.Result,
			Notif:  v.Notif,
		}
	}
	return serialized
}
