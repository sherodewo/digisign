package routes

import (
	"gopkg.in/go-playground/validator.v9"
	"kpdigisign/app/request"
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

type Validator struct {
	validator *validator.Validate
}

func customValidation(vd validator.StructLevel) {
	losValidator := vd.Current().Interface().(request.LosRequest)
	if losValidator.KonsumenType == "NEW" {
		if losValidator.AsliRiRefVerifikasi == nil {
			vd.ReportError(losValidator.AsliRiRefVerifikasi, "asliri_ref_verifikasi", "AsliRiRefVerifikasi", "Cannot be null or empty", "")
		}
		if losValidator.AsliRiNama == nil {
			vd.ReportError(losValidator.AsliRiNama, "asliri_nama", "AsliRiNama", "Cannot be null or empty", "")
		}
		if losValidator.AsliRiAlamat == nil {
			vd.ReportError(losValidator.AsliRiAlamat, "asliri_alamat", "AsliRiAlamat", "Cannot be null or empty", "")
		}
		if losValidator.AsliRiTempatLahir == nil {
			vd.ReportError(losValidator.AsliRiTempatLahir, "asliri_tempat_lahir", "AsliRiTempatLahir", "Cannot be null or empty", "")
		}
		if losValidator.AsliRiTanggalLahir == nil {
			vd.ReportError(losValidator.AsliRiTanggalLahir, "asliri_tanggal_lahir", "AsliRiTanggalLahir", "Cannot be null or empty", "")
		}
		if losValidator.AsliRiRegNumber == nil {
			vd.ReportError(losValidator.AsliRiRegNumber, "asliri_reg_number", "AsliRiRegNumber", "Cannot be null or empty", "")
		}
		if losValidator.ScoreSelfie == nil {
			vd.ReportError(losValidator.ScoreSelfie, "score_selfie", "ScoreSelfie", "Cannot be null or empty", "")
		}
		if losValidator.Vnik == nil {
			vd.ReportError(losValidator.Vnik, "vnik", "Vnik", "Cannot be null or empty", "")
		}
		if losValidator.Vnama == nil {
			vd.ReportError(losValidator.Vnama, "vnama", "Vnama", "Cannot be null or empty", "")
		}
		if losValidator.VtanggalLahir == nil {
			vd.ReportError(losValidator.VtanggalLahir, "vtanggal_lahir", "VtanggalLahir", "Cannot be null or empty", "")
		}
		if losValidator.VtempatLahir == nil {
			vd.ReportError(losValidator.Vnama, "vtempat_lahir", "VtempatLahir", "Cannot be null or empty", "")
		}
	}
}
func (v *Validator) Validate(i interface{}) error {
	v.validator.RegisterStructValidation(customValidation, request.LosRequest{})

	return v.validator.Struct(i)
}
