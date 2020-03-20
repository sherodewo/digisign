package activation

type Dto struct {
	EmailUser string `json:"email_user" validate:"required"`
}
