package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
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
func (us *LosService) Create(u request.LosRequest) (los models.Los, err error) {

	los.ProspectID = u.ProspectID
	los.UserID = u.UserID
	los.Alamat = u.Alamat
	los.JenisKelamin= u.JenisKelamin
	los.Kecamatan = u.Kecamatan
	los.Kelurahan= u.Kelurahan
	los.KodePos = u.KodePos
	los.Kota = u.Kota
	los.Nama= u.Nama
	los.NoTelepon = u.NoTelepon
	los.TanggalLahir = u.TanggalLahir
	los.Provinsi= u.Provinsi
	los.Nik = u.Nik
	los.TempatLahir = u.TempatLahir
	los.Email = u.Email
	los.Npwp = u.Npwp
	los.RegNumber = u.RegNumber
	los.KonsumenType= u.KonsumenType
	los.AsliRiRegNumber = u.AsliRiRegNumber
	los.AsliRiRefVerifikasi = u.AsliRiRefVerifikasi
	los.EmailBm = u.EmailBm
	los.BranchID = u.BranchID
	los.AsliRiSelfieSimilarity = u.AsliRiSelfieSimilarity
	los.AsliRiAlamat = u.AsliRiAlamat
	los.AsliRiTempatLahir = u.AsliRiTempatLahir
	los.AsliRiTanggalLahir = u.AsliRiTanggalLahir
	los.AsliRiNama = u.AsliRiNama

	err = us.db.Create(&los).Error
	return los, err
}

func (us *LosService) Update(id string, um request.LosRequest) (*models.Los, error) {
	data, err := us.FindById(id)
	if err != nil {
		return &data, err
	}
	data.Email = um.Email

	if err := us.db.Save(&data).Error; err != nil {
		log.Debug("ERROR", err)
		return &data, err
	}
	return &data, nil
}

func (us *LosService) FindById(id string) (c models.Los, err error) {
	if err := us.db.Where("id = ?", id).First(&c).Error; err != nil {
		log.Debug("ERROR", err)
		return c, err
	}
	return c, err
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

func (us *LosService) SaveResult(result string, notif string,jsonResponse string) (losResult models.DigisignResult, err error) {

	losResult.Result = result
	losResult.Notif = result
	losResult.JsonResponse = jsonResponse

	if err := us.db.Create(&losResult).Error; err != nil {
		return losResult, err
	}
	return losResult, err
	}