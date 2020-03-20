package download_document

type Dto struct {
	DocumentID string `json:"document_id" validate:"required"`
}
