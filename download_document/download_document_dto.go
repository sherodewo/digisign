package download_document

type Dto struct {
	UserID     string `json:"user_id" validate:"required"`
	DocumentID string `json:"document_id" validate:"required"`
}
