package request

type Register struct {
	Address    string `json:"address" validate:"required,max=50"`
	Gender     string `json:"gender" validate:"validate:required,max=10"`
	Kecamatan  string `json:"kecamatan" validate:"required"`
	Kelurahan  string `json:"kelurahan" validate:"required"`
	Zipcode    string `json:"zipcode" validate:"required,max=10"`
	City       string `json:"city" validate:"required"`
	Name       string `json:"name" validate:"required,max=80"`
	Phone      string `json:"mobile_phone" validate:"required,max=15"`
	BirthDate  string `json:"birth_date" validate:"required,max=10"`
	Provinci   string `json:"provinci" validate:"required"`
	IDKtp      string `json:"id_ktp" validate:"len=16"`
	BirthPlace string `json:"birth_place" validate:"required,max=30"`
	Email      string `json:"email" validate:"required,max=80"`
	NPWP       string `json:"npwp"`
	PhotoKTP   string `json:"photo_ktp"`
	Selfie     string `json:"selfie"`
	Signature  string `json:"signature"`
	PhotoNPWP  string `json:"photo_npwp"`
}

type RegisterRequest struct {
	JsonFile struct {
		UserID     string `json:"user_id" validate:"email,max=80"`
		Address    string `json:"alamat" validate:"required,max=50"`
		Gender     string `json:"jenis_kelamin" validate:"validate:required,max=10"`
		Kecamatan  string `json:"kecamatan" validate:"required"`
		Kelurahan  string `json:"kelurahan" validate:"required"`
		Zipcode    string `json:"kode-pos" validate:"required,max=10"`
		City       string `json:"kota" validate:"required"`
		Name       string `json:"nama" validate:"required,max=80"`
		Phone      string `json:"tlp" validate:"required,max=15"`
		TglLahir   string `json:"tgl_lahir" validate:"required,max=10"`
		Provinci   string `json:"provinci" validate:"required"`
		IDKtp      string `json:"idktp" validate:"len=16"`
		BirthPlace string `json:"tmp_lahir" validate:"required,max=30"`
		Email      string `json:"email" validate:"required,max=80"`
		NPWP       string `json:"npwp"`
		Redirect   bool   `json:"redirect"`
	} `json:"JSONFile"`
}

type ActivationRequest struct {
	JsonFile struct {
		UserID string `json:"user_id" validate:"email,max=80"`
		Email  string `json:"email_user" validate:"required,max=80"`
	} `json:"JSONFile"`
}

type SendDocRequest struct {
	JsonFile struct {
		UserID         string    `json:"user_id" validate:"email,max=80"`
		DocumentID     string    `json:"document_id" validate:"required,max=20"`
		Payment        string    `json:"payment" validate:"max=1"`
		Redirect       bool      `json:"redirect"`
		Branch         string    `json:"branch"`
		SequenceOption bool      `json:"sequence_option"`
		SendTo         []SendTo  `json:"sent-to"`
		ReqSign        []ReqSign `json:"req-sign"`
	} `json:"JSONFile"`
}

type SendTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReqSign struct {
	Name       string `json:"name" validate:"required,max=80"`
	Email      string `json:"email" validate:"required,max=80"`
	User       string `json:"user" validate:"required,max=30"`
	Llx        string `json:"llx" validate:"max=30"`
	Lly        string `json:"lly" validate:"max=30"`
	Urx        string `json:"urx" validate:"max=30"`
	Ury        string `json:"ury" validate:"max=30"`
	Page       string `json:"page" validate:"max=3"`
	Visible    string `json:"visible" validate:"max=1"`
	SigningSeq int    `json:"signing_seq"`
}

type SignDocRequest struct {
	JsonFile JsonFileSign `json:"JSONFile"`
}

type DownloadRequest struct {
	JSONFile DownloadDto `json:"JSONFile"`
}

type UploadMediaRequest struct {
	Type        string `json:"type"`
	ReferenceNo string `json:"reference_no"`
	File        string `json:"file"`
}

type JsonFileSign struct {
	UserID     string `json:"userid" validate:"email,max=80"`
	DocumentID string `json:"document_id" validate:"required,max=20"`
	Email      string `json:"email_user" validate:"required,max=80"`
	ViewOnly   bool   `json:"view_only"`
}

type SignDocDto struct {
	ProspectID string `json:"prospect_id" validate:"required,max=50"`
	UserID     string `json:"user_id" validate:"email,max=80"`
	DocumentID string `json:"document_id" validate:"required,max=20"`
	Email      string `json:"email_user" validate:"required,max=80"`
	ViewOnly   bool   `json:"view_only"`
}

type DownloadDto struct {
	UserID     string `json:"user_id" validate:"email,max=80"`
	DocumentID string `json:"document_id" validate:"required,max=20"`
}
