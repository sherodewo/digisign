package send_document

type Dto struct {
	UserID         string    `form:"user_id" json:"user_id" validate:"required"`
	DocumentID     string    `form:"document_id" json:"document_id" validate:"required"`
	Payment        string    `form:"payment" json:"payment" validate:"required"`
	SendTo         []SendTo  `json:"send_to" validate:"required,dive,required"`
	ReqSign        []ReqSign `json:"req_sign" validate:"required,dive,required"`
	Redirect       bool      `json:"redirect"`
	SequenceOption bool      `json:"sequence_option"`
	SigningSeq     int       `json:"signing_seq"`
	File           string    `json:"file" validate:"required"`
}
type SendTo struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type ReqSign struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required"`
	AksiTtd string `json:"aksi_ttd"`
	Kuser   string `json:"kuser"`
	User    string `json:"user" validate:"required"`
	Page    string `json:"page"`
	Llx     string `json:"llx"`
	Lly     string `json:"lly"`
	Urx     string `json:"urx"`
	Ury     string `json:"ury"`
	Visible string `json:"visible"`
}
