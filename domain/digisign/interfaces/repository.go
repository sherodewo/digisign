package interfaces

import "los-int-digisign/model/entity"

type Repository interface {
	GetSendDocData(prospectID string) (data entity.SendDocData, err error)
	GetCustomerPersonalByEmailAndNik(email, nik string) (data entity.CustomerPersonal, err error)
}
