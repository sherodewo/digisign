package request

type LosRequest struct {
	ProspectID             string `json:"prospect_id"`
	UserID                 string `json:"user_id"`
	Alamat                 string `json:"alamat"`
	JenisKelamin           string `json:"jenis_kelamin"`
	Kecamatan              string `json:"kecamatan" `
	Kelurahan              string `json:"kelurahan "`
	KodePos                string `json:"kode_pos"`
	Kota                   string `json:"kota"`
	Nama                   string `json:"nama"`
	NoTelepon              string `json:"no_telepon"`
	TanggalLahir           string `json:"tanggal_lahir"`
	Provinsi               string `json:"provinsi"`
	Nik                    string `json:"nik" `
	TempatLahir            string `json:"tempat_lahir"`
	Email                  string `json:"email"`
	Npwp                   string `json:"npwp"`
	RegNumber              string `json:"reg_number"`
	KonsumenType           string `json:"konsumen_type"`
	AsliRiRegNumber        string `json:"asliri_reg_number"`
	AsliRiRefVerifikasi    int    `json:"asliri_ref_verifikasi"`
	AsliRiNama             bool   `json:"asliri_nama"`
	AsliRiTempatLahir      bool   `json:"asliri_tempat_lahir"`
	AsliRiTanggalLahir     bool   `json:"asliri_tanggal_lahir"`
	AsliRiAlamat           string `json:"asliri_alamat"`
	AsliRiSelfieSimilarity string `json:"asliri_selfie_similarity"`
	BranchID               string `json:"branch_id"`
	EmailBm                string `json:"email_bm"`
}

type LosSendDocumentRequest struct {
	Userid     string `form:"userid"`
	Documentid string `form:"documentid"`
	Payment    string `form:"payment"`
	SendTo     string `form:"sendTo"`
	ReqSign    string `form:"reqSign"`
}

type LosDownloadDocumentRequest struct {
	UserID     string `json:"user_id"`
	DocumentID string `json:"document_id"`
}

/*type LosSendTo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LosReqSign struct {
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
}*/
