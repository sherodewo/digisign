package request

type LosRequest struct {
	ProspectID          string  `json:"prospect_id" validate:"required"`
	Alamat              string  `json:"alamat" validate:"required"`
	JenisKelamin        string  `json:"jenis_kelamin" validate:"required"`
	Kecamatan           string  `json:"kecamatan" validate:"required"`
	Kelurahan           string  `json:"kelurahan" validate:"required"`
	KodePos             string  `json:"kode_pos" validate:"required"`
	Kota                string  `json:"kota" validate:"required"`
	Nama                string  `json:"nama" validate:"required"`
	NoTelepon           string  `json:"no_telepon" validate:"required"`
	TanggalLahir        string  `json:"tanggal_lahir" validate:"required"`
	Provinsi            string  `json:"provinsi" validate:"required"`
	Nik                 string  `json:"nik" validate:"required"`
	TempatLahir         string  `json:"tempat_lahir" validate:"required"`
	Email               string  `json:"email" validate:"required"`
	Npwp                string  `json:"npwp"`
	RegNumber           string  `json:"reg_number" validate:"required"`
	KonsumenType        string  `json:"konsumen_type" validate:"required"`
	AsliRiRegNumber     *string `json:"asliri_reg_number"`
	AsliRiRefVerifikasi *string `json:"asliri_ref_verifikasi"`
	AsliRiNama          *bool   `json:"asliri_nama"`
	AsliRiTempatLahir   *bool   `json:"asliri_tempat_lahir"`
	AsliRiTanggalLahir  *bool   `json:"asliri_tanggal_lahir"`
	AsliRiAlamat        *string `json:"asliri_alamat"`
	ScoreSelfie         *string `json:"score_selfie"`
	Vnik                *string `json:"vnik"`
	Vnama               *string `json:"vnama"`
	VtanggalLahir       *string `json:"vtanggal_lahir"`
	VtempatLahir        *string `json:"vtempat_lahir"`
	BranchID            string  `json:"branch_id" validate:"required"`
	EmailBm             string  `json:"email_bm" validate:"required"`
}

type LosSendDocumentRequest struct {
	DocumentID string `form:"documentId" validate:"required"`
	Payment    string `form:"payment" validate:"required"`
	SendTo     string `form:"sendTo" validate:"required"`
	ReqSign    string `form:"reqSign" validate:"required"`
}

type LosDownloadDocumentRequest struct {
	DocumentID string `json:"document_id" validate:"required"`
}

type LosActivationRequest struct {
	EmailUser string `json:"email_user" validate:"required"`
}

type LosSignDocumentRequest struct {
	EmailUser  string `json:"email_user" validate:"required"`
	DocumentID string `json:"document_id" validate:"required"`
}
