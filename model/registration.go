package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type Registration struct {
	ID                   string             `gorm:"column:id;primary_key:true"`
	RegistrationResult   RegistrationResult `gorm:"forigenkey:RegistrationResultID"`
	RegistrationResultID string             `gorm:"column:registration_result_id"`
	ProspectID           string             `gorm:"type:varchar(100);column:prospect_id"`
	UserID               string             `gorm:"type:varchar(100);column:user_id"`
	Alamat               string             `gorm:"column:alamat"`
	JenisKelamin         string             `gorm:"type:varchar(12);column:jenis_kelamin"`
	Kecamatan            string             `gorm:"type:varchar(50);column:kecamatan"`
	Kelurahan            string             `gorm:"type:varchar(50);column:kelurahan"`
	KodePos              string             `gorm:"type:varchar(20);column:kode_pos"`
	Kota                 string             `gorm:"type:varchar(50);column:kota"`
	Nama                 string             `gorm:"type:varchar(150);column:nama"`
	NoTelepon            string             `gorm:"type:varchar(20);column:no_telepon"`
	TanggalLahir         string             `gorm:"type:varchar(50);column:tanggal_lahir"`
	Provinsi             string             `gorm:"type:varchar(50);column:provinsi"`
	Nik                  string             `gorm:"type:varchar(100);column:nik"`
	TempatLahir          string             `gorm:"type:varchar(50);column:tempat_lahir"`
	Email                string             `gorm:"type:varchar(50);column:email"`
	Npwp                 string             `gorm:"type:varchar(50);column:npwp"`
	RegNumber            string             `gorm:"type:varchar(100);column:reg_number"`
	KonsumenType         string             `gorm:"type:varchar(20);column:konsumen_type"`
	AsliRiRegNumber      string             `gorm:"type:varchar(150);column:asliri_reg_number"`
	AsliRiRefVerifikasi  string             `gorm:"column:asliri_ref_verifikasi"`
	AsliRiNama           bool               `gorm:"column:asliri_nama"`
	AsliRiTempatLahir    bool               `gorm:"column:asliri_tempat_lahir"`
	AsliRiTanggalLahir   bool               `gorm:"column:asliri_tanggal_lahir"`
	AsliRiAlamat         string             `gorm:"type:varchar(255);column:asliri_alamat"`
	ScoreSelfie          string             `gorm:"type:varchar(50);column:score_selfie"`
	Vnik                 string             `gorm:"type:varchar(2);column:vnik"`
	Vnama                string             `gorm:"type:varchar(2);column:vnama"`
	VtanggalLahir        string             `gorm:"type:varchar(2);column:vtanggal_lahir"`
	VtempatLahir         string             `gorm:"type:varchar(2);column:vtempat_lahir"`
	BranchID             string             `gorm:"type:varchar(50);column:branch_id"`
	EmailBm              string             `gorm:"type:varchar(50);column:emai_bm"`
	Redirect             bool               `gorm:"column:redirect"`
	CreatedAt            time.Time          `gorm:"column:created_at;"`
	UpdatedAt            time.Time          `gorm:"column:updated_at;"`
}

func (m *Registration) TableName() string {
	return "registrations"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (m *Registration) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
