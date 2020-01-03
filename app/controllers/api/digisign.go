package api

import (
	"github.com/labstack/echo"
	"kpdigisign/app/client"
	"kpdigisign/app/helpers"
	"kpdigisign/app/repository"
	"kpdigisign/app/request"
	"kpdigisign/app/response"
)

type DigisignController struct {
	LosRepository repository.LosRequestRepository
}

func (d *DigisignController) Register(c echo.Context) error {

	losRequest := request.LosRequest{}
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

	register := client.NewDigisignRegistrationRequest()
	resp, err := register.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)

	if err != nil {
		return response.BadRequest(c, "Bad Request", err.Error())
	}
	return response.SingleData(c, "Success Execute request", resp.String())
}

func (d *DigisignController) Check(c echo.Context) error  {
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
	resLos,_ := checkLosRequest.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
	losRespone := response.NewLosRespone()

	if err := losRespone.Bind(resLos.Body()); err != nil {
		return response.BadRequest(c,"Bad Request", err)
	}

	//Check Konsumen Type
	if losRequest.KonsumenType == "NEW" {
		register := client.NewDigisignRegistrationRequest()
		resp, err := register.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
		if err != nil {
			return response.BadRequest(c, "Bad Request", err.Error())
		}
		_, _ = d.LosRepository.SaveResult(losRespone.Result, losRespone.Info, losRespone.EmailRegistered, losRespone.Name,
			losRespone.Birthplace, losRespone.Birthdate, losRespone.Address, losRespone.SelfieMatch)
		return response.SingleData(c, "Success Execute request", resp.String())

	}else {
		register := client.NewDigisignRegistrationRequestRoAo()
		resp, err := register.DigisignRegistration(bufKtp, bufSelfie, bufNpwp, bufTtd, losRequest)
		if err != nil {
			return response.BadRequest(c, "Bad Request", err.Error())
		}
		_, _ = d.LosRepository.SaveResult(losRespone.Result, losRespone.Info, losRespone.EmailRegistered, losRespone.Name,
			losRespone.Birthplace, losRespone.Birthdate, losRespone.Address, losRespone.SelfieMatch)
		return response.SingleData(c, "Success Execute request", resp.String())
	}


	return nil
}