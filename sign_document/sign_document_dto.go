package sign_document

type Dto struct {
	EmailUser  string `json:"email_user" validate:"required"`
	DocumentID string `json:"document_id" validate:"required"`
}
