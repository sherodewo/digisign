package registration

import (
	"los-int-digisign/infrastructure/config/digisign"
	"los-int-digisign/model"
)

type service struct {
	registrationRepository Repository
	registrationMapper     *Mapper
}

func NewRegistrationService(repository Repository, mapper *Mapper) *service {
	return &service{
		registrationRepository: repository,
		registrationMapper:     mapper,
	}
}

func (s *service) FindAllRegistrations() (*[]Mapper, error) {
	data, err := s.registrationRepository.FindAll()
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "registration",
			"app.func": "FindAllRegistrations",
			"action":   "read",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.registrationMapper.MapList(data), nil
	}
}

func (s *service) FindRegistrationById(id string) (*Mapper, error) {
	data, err := s.registrationRepository.FindById(id)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "registration",
			"app.func": "FindRegistrationById",
			"action":   "read",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message": err.Error(),
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.registrationMapper.Map(data), nil
	}
}

func (s *service) SaveRegistration(dto Dto, result string, notif string, reftrx string,
	jsonResponse string, kodeUser string, jsonRequest string) (*Mapper, error) {
	entity := model.Registration{}
	entity.ProspectID = dto.ProspectID
	entity.UserID = dto.UserID
	entity.Alamat = dto.Alamat
	entity.JenisKelamin = dto.JenisKelamin
	entity.Kecamatan = dto.Kecamatan
	entity.Kelurahan = dto.Kelurahan
	entity.KodePos = dto.KodePos
	entity.Kota = dto.Kota
	entity.Nama = dto.Nama
	entity.NoTelepon = dto.NoTelepon
	entity.TanggalLahir = dto.TanggalLahir
	entity.Provinsi = dto.Provinsi
	entity.Nik = dto.Nik
	entity.TempatLahir = dto.TempatLahir
	entity.Email = dto.Email
	entity.Npwp = dto.Npwp
	entity.KonsumenType = dto.KonsumenType
	entity.Redirect = dto.Redirect
	if entity.KonsumenType == "NEW" {
		entity.AsliRiRefVerifikasi = *dto.AsliRiRefVerifikasi
		entity.ScoreSelfie = *dto.ScoreSelfie
		entity.AsliRiAlamat = *dto.AsliRiAlamat
		entity.AsliRiTempatLahir = *dto.AsliRiTempatLahir
		entity.AsliRiTanggalLahir = *dto.AsliRiTanggalLahir
		entity.AsliRiNama = *dto.AsliRiNama
		entity.Vnama = *dto.Vnama
		entity.Vnik = *dto.Vnik
		entity.VtanggalLahir = *dto.VtanggalLahir
		entity.VtempatLahir = *dto.VtempatLahir
	}
	entity.RegistrationResult.RefTrx = reftrx
	entity.RegistrationResult.Notif = notif
	entity.RegistrationResult.Result = result
	entity.RegistrationResult.KodeUser = kodeUser
	entity.RegistrationResult.JsonResponse = jsonResponse
	entity.RegistrationResult.JsonRequest = jsonRequest

	data, err := s.registrationRepository.Save(entity)
	if err != nil {
		tags := map[string]string{
			"app.pkg":  "registration",
			"app.func": "SaveRegistration",
			"action":   "create",
			"db.name":  "di****gn",
		}
		extra := map[string]interface{}{
			"message":     err.Error(),
			"prospect_id": entity.ProspectID,
			"user_id":     entity.UserID,
		}
		digisign.SendToSentry(tags, extra, "DATABASE")
		return nil, err
	} else {
		return s.registrationMapper.Map(data), nil
	}
}
