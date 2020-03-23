package registration

type Dto struct {
	UserID              string  `json:"user_id" form:"user_id" validate:"required"`
	ProspectID          string  `json:"prospect_id" form:"prospect_id" validate:"required"`
	Alamat              string  `json:"alamat" form:"alamat" validate:"required"`
	JenisKelamin        string  `json:"jenis_kelamin" form:"jenis_kelamin" validate:"required"`
	Kecamatan           string  `json:"kecamatan" form:"kecamatan" validate:"required"`
	Kelurahan           string  `json:"kelurahan" form:"kelurahan" validate:"required"`
	KodePos             string  `json:"kode_pos" form:"kode_pos" validate:"required"`
	Kota                string  `json:"kota" form:"kota" validate:"required"`
	Nama                string  `json:"nama" form:"nama" validate:"required"`
	NoTelepon           string  `json:"no_telepon" form:"no_telepon" validate:"required"`
	TanggalLahir        string  `json:"tanggal_lahir" form:"tanggal_lahir" validate:"required"`
	Provinsi            string  `json:"provinsi" form:"provinsi" validate:"required"`
	Nik                 string  `json:"nik" form:"nik" validate:"required"`
	TempatLahir         string  `json:"tempat_lahir" form:"tempat_lahir" validate:"required"`
	Email               string  `json:"email" form:"email" validate:"required"`
	Npwp                string  `json:"npwp" form:"npwp"`
	//RegNumber           string  `json:"reg_number" form:"reg_number" validate:"required"`
	KonsumenType        string  `json:"konsumen_type" form:"konsumen_type" validate:"required"`
	//AsliRiRegNumber     *string `json:"asliri_reg_number" form:"asliri_reg_number"`
	AsliRiRefVerifikasi *string `json:"asliri_ref_verifikasi" form:"asliri_ref_verifikasi"`
	AsliRiNama          *bool   `json:"asliri_nama" form:"asliri_nama"`
	AsliRiTempatLahir   *bool   `json:"asliri_tempat_lahir" form:"asliri_tempat_lahir"`
	AsliRiTanggalLahir  *bool   `json:"asliri_tanggal_lahir" form:"asliri_tanggal_lahir"`
	AsliRiAlamat        *string `json:"asliri_alamat" form:"asliri_alamat"`
	ScoreSelfie         *string `json:"score_selfie" form:"score_selfie"`
	Vnik                *string `json:"vnik" form:"vnik"`
	Vnama               *string `json:"vnama" form:"vnama"`
	VtanggalLahir       *string `json:"vtanggal_lahir" form:"vtanggal_lahir"`
	VtempatLahir        *string `json:"vtempat_lahir" form:"vtempat_lahir"`
	//BranchID            string  `json:"branch_id" form:"branch_id" validate:"required"`
	//EmailBm             string  `json:"email_bm" form:"email_bm" validate:"required"`
	FotoKtp             string  `json:"foto_ktp" form:"foto_ktp" validate:"required"`
	FotoSelfie          string  `json:"foto_selfie" form:"foto_selfie" validate:"required"`
	FotoNpwp            string  `json:"foto_npwp" form:"foto_npwp"`
	FotoTandaTangan     string  `json:"foto_tanda_tangan" form:"foto_tanda_tangan"`
	Redirect            bool    `json:"redirect" validate:"required"`
}
