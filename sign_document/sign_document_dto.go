package sign_document

type Dto struct {
	UserID     string `json:"user_id" validate:"required"`
	EmailUser  string `json:"email_user" validate:"required"`
	DocumentID string `json:"document_id" validate:"required"`
}
