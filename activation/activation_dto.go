package activation

type Dto struct {
	UserID     string `json:"user_id" validate:"required"`
	ProspectID string `json:"prospect_id" form:"prospect_id" validate:"required"`
	EmailUser  string `json:"email_user" validate:"required"`
}
