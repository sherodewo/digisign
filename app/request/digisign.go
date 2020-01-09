package request

type DigisignRegistrationRequest struct {
	JsonFile JsonFile `json:"JSONFile"`
}

type JsonFile struct {
	UserID                 string `json:"userid"`
	Alamat                 string `json:"alamat"`
	JenisKelamin           string `json:"jenis_kelamin"`
	Kecamatan              string `json:"kecamatan"`
	Kelurahan              string `json:"kelurahan"`
	KodePos                string `json:"kode-pos"`
	Kota                   string `json:"kota"`
	Nama                   string `json:"nama"`
	NoTelepon              string `json:"no_telepon"`
	TanggalLahir           string `json:"tgl_lahir"`
	Provinsi               string `json:"provinci"`
	Nik                    string `json:"idktp"`
	TempatLahir            string `json:"tmp_lahir"`
	Email                  string `json:"email"`
	Npwp                   string `json:"npwp"`
	RegNumber              string `json:"reg_number"`
	AsliRiRegNumber        string `json:"asliri_reg_number,omitempty"`
	AsliRiRefVerifikasi    int    `json:"asliri_ref_verifikasi,omitempty"`
	AsliRiNama             bool   `json:"asliri_nama,omitempty"`
	AsliRiTempatLahir      bool   `json:"asliri_tempat_lahir,omitempty"`
	AsliRiTanggalLahir     bool   `json:"asliri_tanggal_lahir,omitempty"`
	AsliRiAlamat           string `json:"asliri_alamat,omitempty"`
	AsliRiSelfieSimilarity string `json:"asliri_selfie_similarity,omitempty"`
	BranchID               string `json:"branch_id"`
	EmailBm                string `json:"email_bm"`
}

type SendDocumentRequest struct {
	JsonFile JsonFileDoc `json:"JSONFile"`
}

type JsonFileDoc struct {
	UserID     string                 `json:"userid"`
	DocumentID string                 `json:"document_id"`
	Payment    string                 `json:"payment"`
	SendTo     interface{} `json:"send-to"`
	ReqSign    interface{} `json:"req-sign"`
}

type SendTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ReqSign struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	AksiTtd string `json:"aksi_ttd"`
	Kuser   string `json:"kuser"`
	User    string `json:"user"`
	Page    string `json:"page"`
	Llx     string `json:"llx"`
	Lly     string `json:"lly"`
	Urx     string `json:"urx"`
	Ury     string `json:"ury"`
	Visible string `json:"visible"`
}
