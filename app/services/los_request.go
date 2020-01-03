package service

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"kpdigisign/app/models"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
)

// UserService :
type LosRequestService struct {
	db *gorm.DB
}

// NewUserService :
func NewLosRequestService(db *gorm.DB) repository.LosRequestRepository {
	return &LosRequestService{
		db,
	}
}

// GetByID :
func (us *LosRequestService) GetByID(id string) (*models.LosRequest, error) {
	var m models.LosRequest
	if err := us.db.Where(&models.LosRequest{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}
	return &m, nil
}

// GetByEmail :
func (us *LosRequestService) GetByEmail(e string) (*models.LosRequest, error) {
	var m models.LosRequest
	if err := us.db.Where(&models.LosRequest{Email: e}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, err
	}
	return &m, nil
}

// Create :
func (us *LosRequestService) Create(u request.LosRequest) (los models.LosRequest, err error) {

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

func (us *LosRequestService) Update(id string, um request.LosRequest) (*models.LosRequest, error) {
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

func (us *LosRequestService) FindById(id string) (c models.LosRequest, err error) {
	if err := us.db.Where("id = ?", id).First(&c).Error; err != nil {
		log.Debug("ERROR", err)
		return c, err
	}
	return c, err
}

func (us *LosRequestService) FindAll() (user []models.LosRequest, err error) {
	err = us.db.Find(&user).Error
	return user, err
}

func (us *LosRequestService) Destroy(id string) error {
	var m models.LosRequest
	if err := us.db.Where(&models.LosRequest{ID: id}).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return err
		}
		return err
	}
	if err := us.db.Delete(&models.LosRequest{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func (us *LosRequestService) SaveResult(result string, info string, emailRegistered string, name bool, birthplace bool,
	birthdate bool, address string, selfieMatch bool) (losResult models.DigisignRegistrationResult, err error) {
	losResult.Name = name
	losResult.Address = address
	losResult.Result = result
	losResult.Info = info
	losResult.BirthDate = birthdate
	losResult.BirthPlace = birthplace
	losResult.EmailRegistered = emailRegistered
	losResult.SelfieMatch = selfieMatch
	if err := us.db.Create(&losResult).Error; err != nil {
		return losResult, err
	}
	return losResult, err
	}