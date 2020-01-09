package request

type DownloadRequest struct {
	JsonFile JsonDownloadFile `json:"JSONFile"`
}
type JsonDownloadFile struct {
	UserID     string                 `json:"userid"`
	DocumentID string                 `json:"document_id"`
}
