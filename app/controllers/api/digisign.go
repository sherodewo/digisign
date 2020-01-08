package api

import (
	"github.com/labstack/echo"
	"kpdigisign/app/client"
	"kpdigisign/app/helpers"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
	"kpdigisign/app/response"
	"kpdigisign/app/response/mapper"
)

type DigisignController struct {
	LosRepository repository.LosRepository
}

func (d *DigisignController) Register(c echo.Context) error {
	resultMapper := mapper.NewDigisignResultMapper()
	losRequest := request.LosRequest{}
	if err := c.Bind(&losRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	var bufNpwp, bufTtd []byte

	//Check KTP
	fileKtp, err := c.FormFile("foto_ktp")
	if fileKtp == nil {
		return response.BadRequest(c, "NOT FOUND KTP", nil)
	}
	bufKtp, err := helpers.GetImageByte("foto_ktp", c)
	//Check Selfie
	fileSelfie, err := c.FormFile("foto_selfie")
	if fileSelfie == nil {
		return response.BadRequest(c, "NOT FOUND Selfie", nil)
	}
	bufSelfie, err := helpers.GetImageByte("foto_selfie", c)
	//Get NPWP
	bufNpwp, err = helpers.GetImageByte("foto_npwp", c)
	//Get TTD
	bufTtd, err = helpers.GetImageByte("tanda_tangan", c)

	_, err = d.LosRepository.Create(losRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	register := client.NewDigisignRegistrationRequest()
	resp, err := register.DigisignRegistration(losRequest.KonsumenType, bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)

	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}

	respDigisignRegister := response.NewDigisignResponse()
	if err := respDigisignRegister.Bind(resp.Body()); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	resultData, err := d.LosRepository.SaveResult(respDigisignRegister.JsonFile.Result, respDigisignRegister.JsonFile.Notif, resp.String())

	return response.SingleData(c, "Success execute resuest", resultMapper.Map(resultData))
}

/*func (d *DigisignController) Check(c echo.Context) error {
	losRequest := request.LosRequest{}
	//resultMapper := mapper.NewDigisignRegistrationResultMapper()
	if err := c.Bind(&losRequest); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	bufKtp, err := helpers.GetImageByte("foto_ktp", c)
	bufSelfie, err := helpers.GetImageByte("foto_selfie", c)
	bufNpwp, err := helpers.GetImageByte("foto_npwp", c)
	bufTtd, err := helpers.GetImageByte("tanda_tangan", c)

	_, err = d.LosRepository.Create(losRequest)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	//Mapping Los
	checkLosRequest := client.NewDigisignRegistrationRequest()
	resLos, _ := checkLosRequest.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
	losRespone := response.NewLosRespone()

	if err := losRespone.Bind(resLos.Body()); err != nil {
		return response.BadRequest(c, "Bad Request", err)
	}

	//Check Konsumen Type
	if losRequest.KonsumenType == "NEW" {
		register := client.NewDigisignRegistrationRequest()
		resp, err := register.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
		if err != nil {
			return response.BadRequest(c, "Bad Request", err.Error())
		}
		_, _ = d.LosRepository.SaveResult(losRespone.JsonFile.Result, losRespone.JsonFile.Info, losRespone.JsonFile.EmailRegistered, losRespone.JsonFile.Name,
			losRespone.JsonFile.Birthplace, losRespone.JsonFile.Birthdate, losRespone.JsonFile.Address, losRespone.JsonFile.SelfieMatch)
		return response.SingleData(c, "Success Execute request", resp.String())

	} else {
		register := client.NewDigisignRegistrationRequestRoAo()
		resp, err := register.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
		if err != nil {
			return response.BadRequest(c, "Bad Request", err.Error())
		}
		_, _ = d.LosRepository.SaveResult(losRespone.JsonFile.Result, losRespone.JsonFile.Info, losRespone.JsonFile.EmailRegistered, losRespone.Name,
			losRespone.JsonFile.Birthplace, losRespone.JsonFile.Birthdate, losRespone.JsonFile.Address, losRespone.JsonFile.SelfieMatch)
		return response.SingleData(c, "Success Execute request", resp.String())
	}

	return nil
}*/
