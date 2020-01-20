package request

type DigisignRegistrationRequest struct {
	JsonFile JsonFile `json:"JSONFile"`
}

type JsonFile struct {
	UserID              string `json:"userid"`
	Alamat              string `json:"alamat"`
	JenisKelamin        string `json:"jenis_kelamin"`
	Kecamatan           string `json:"kecamatan"`
	Kelurahan           string `json:"kelurahan"`
	KodePos             string `json:"kode-pos"`
	Kota                string `json:"kota"`
	Nama                string `json:"nama"`
	NoTelepon           string `json:"tlp"`
	TanggalLahir        string `json:"tgl_lahir"`
	Provinsi            string `json:"provinci"`
	Nik                 string `json:"idktp"`
	TempatLahir         string `json:"tmp_lahir"`
	Email               string `json:"email"`
	Npwp                string `json:"npwp"`
	RegNumber           string `json:"reg_number"`
	Redirect            bool   `json:"redirect"`
	AsliRiRefVerifikasi string `json:"ref_verifikasi,omitempty"`
	DataVerifikasi      string `json:"data_verifikasi,omitempty"`
	ScoreSelfie         string `json:"score_selfie,omitempty"`
	Vnik                string `json:"vnik,omitempty"`
	Vnama               string `json:"vnama,omitempty"`
	VtanggalLahir       string `json:"vtgl_lahir,omitempty"`
	VtempatLahir        string `json:"vtmp_lahir,omitempty"`
}

type SendDocumentRequest struct {
	JsonFile JsonFileDoc `json:"JSONFile"`
}

type JsonFileDoc struct {
	UserID         string      `json:"userid"`
	DocumentID     string      `json:"document_id"`
	Payment        string      `json:"payment"`
	SendTo         interface{} `json:"send-to"`
	ReqSign        interface{} `json:"req-sign"`
	Redirect       bool        `json:"redirect"`
	SequenceOption bool        `json:"sequence_option"`
}

type ActivationRequest struct {
	JSONFile struct {
		UserID    string `json:"userid"`
		EmailUser string `json:"email_user"`
	}
}

type SignDocumentRequest struct {
	JSONFile struct {
		UserID     string `json:"userid"`
		DocumentID string `json:"document_id"`
		EmailUser  string `json:"email_user"`
		ViewOnly   bool `json:"view_only"`
	}
}
