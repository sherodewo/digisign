package request

import (
	"github.com/labstack/echo"
)

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

func (cr *LosRequest) Bind(c echo.Context) (*LosRequest, error) {
	if err := c.Bind(cr); err != nil {
		return nil, err
	}
	return cr, nil
}
