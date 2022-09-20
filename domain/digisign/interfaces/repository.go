package interfaces

import "los-int-digisign/model/entity"

type Repository interface {
	GetSendDocData(prospectID string) (data entity.SendDocData, err error)
	GetCustomerPersonalByEmailAndNik(email, nik string) (data entity.CustomerPersonal, err error)
	UpdateStatusDigisignActivation(prospectID string) error
	SaveTrx(data []entity.TrxDetail) (err error)
	GetCustomerPersonalByEmail(documentID string) (data entity.TrxDetail, err error)
	GetDigisignDummy(email string, action string) (data entity.DigisignDummy, err error)
}
