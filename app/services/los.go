package service

import (
	"github.com/jinzhu/gorm"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
	"os"
)

// UserService :
type LosService struct {
	db *gorm.DB
}

// NewUserService :
func NewLosRequestService(db *gorm.DB) repository.LosRepository {
	return &LosService{
		db,
	}
}

// GetByID :
func (us *LosService) GetByID(id string) (*models.Los, error) {
	var m models.Los
	if err := us.db.Where(&models.Los{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return &m, nil
}

// GetByEmail :
func (us *LosService) GetByEmail(e string) (*models.Los, error) {
	var m models.Los
	if err := us.db.Where(&models.Los{Email: e}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, err
	}
	return &m, nil
}

// Create :
func (us *LosService) Create(u *request.LosRequest) (los models.Los, err error) {

	los.ProspectID = u.ProspectID
	los.UserID = os.Getenv("DIGISIGN_USER_ID")
	los.Alamat = u.Alamat
	los.JenisKelamin = u.JenisKelamin
	los.Kecamatan = u.Kecamatan
	los.Kelurahan = u.Kelurahan
	los.KodePos = u.KodePos
	los.Kota = u.Kota
	los.Nama = u.Nama
	los.NoTelepon = u.NoTelepon
	los.TanggalLahir = u.TanggalLahir
	los.Provinsi = u.Provinsi
	los.Nik = u.Nik
	los.TempatLahir = u.TempatLahir
	los.Email = u.Email
	los.Npwp = u.Npwp
	los.RegNumber = u.RegNumber
	los.KonsumenType = u.KonsumenType
	los.EmailBm = u.EmailBm
	los.BranchID = u.BranchID
	if los.KonsumenType =="NEW" {
		los.AsliRiRegNumber = *u.AsliRiRegNumber
		los.AsliRiRefVerifikasi = *u.AsliRiRefVerifikasi
		los.ScoreSelfie = *u.ScoreSelfie
		los.AsliRiAlamat = *u.AsliRiAlamat
		los.AsliRiTempatLahir = *u.AsliRiTempatLahir
		los.AsliRiTanggalLahir = *u.AsliRiTanggalLahir
		los.AsliRiNama = *u.AsliRiNama
		los.Vnama = *u.Vnama
		los.Vnik = *u.Vnik
		los.VtanggalLahir = *u.VtanggalLahir
		los.VtempatLahir = *u.VtempatLahir
	}

	err = us.db.Create(&los).Error
	return los, err
}

func (us *LosService) FindAll() (user []models.Los, err error) {
	err = us.db.Find(&user).Error
	return user, err
}

func (us *LosService) Destroy(id string) error {
	var m models.Los
	if err := us.db.Where(&models.Los{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return err
		}
		return err
	}
	if err := us.db.Delete(&models.Los{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
