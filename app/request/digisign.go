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
