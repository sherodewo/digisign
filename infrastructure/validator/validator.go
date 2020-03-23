package validator

import (
	"gopkg.in/go-playground/validator.v9"
	"kpdigisign/registration"
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

type Validator struct {
	validator *validator.Validate
}

func customRegisterValidation(vd validator.StructLevel) {
	registerDto := vd.Current().Interface().(registration.Dto)
	if registerDto.KonsumenType == "NEW" {
		if registerDto.AsliRiRefVerifikasi == nil {
			vd.ReportError(registerDto.AsliRiRefVerifikasi, "asliri_ref_verifikasi", "AsliRiRefVerifikasi", "Cannot be null or empty", "")
		}
		if registerDto.AsliRiNama == nil {
			vd.ReportError(registerDto.AsliRiNama, "asliri_nama", "AsliRiNama", "Cannot be null or empty", "")
		}
		if registerDto.AsliRiAlamat == nil {
			vd.ReportError(registerDto.AsliRiAlamat, "asliri_alamat", "AsliRiAlamat", "Cannot be null or empty", "")
		}
		if registerDto.AsliRiTempatLahir == nil {
			vd.ReportError(registerDto.AsliRiTempatLahir, "asliri_tempat_lahir", "AsliRiTempatLahir", "Cannot be null or empty", "")
		}
		if registerDto.AsliRiTanggalLahir == nil {
			vd.ReportError(registerDto.AsliRiTanggalLahir, "asliri_tanggal_lahir", "AsliRiTanggalLahir", "Cannot be null or empty", "")
		}
		if registerDto.ScoreSelfie == nil {
			vd.ReportError(registerDto.ScoreSelfie, "score_selfie", "ScoreSelfie", "Cannot be null or empty", "")
		}
		if registerDto.Vnik == nil {
			vd.ReportError(registerDto.Vnik, "vnik", "Vnik", "Cannot be null or empty", "")
		}
		if registerDto.Vnama == nil {
			vd.ReportError(registerDto.Vnama, "vnama", "Vnama", "Cannot be null or empty", "")
		}
		if registerDto.VtanggalLahir == nil {
			vd.ReportError(registerDto.VtanggalLahir, "vtanggal_lahir", "VtanggalLahir", "Cannot be null or empty", "")
		}
		if registerDto.VtempatLahir == nil {
			vd.ReportError(registerDto.Vnama, "vtempat_lahir", "VtempatLahir", "Cannot be null or empty", "")
		}
	}
}
func (v *Validator) Validate(i interface{}) error {
	v.validator.RegisterStructValidation(customRegisterValidation, registration.Dto{})

	return v.validator.Struct(i)
}
