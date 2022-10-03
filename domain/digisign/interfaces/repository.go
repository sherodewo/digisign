package interfaces

import "los-int-digisign/model/entity"

type Repository interface {
	GetSendDocData(prospectID string) (data entity.SendDocData, err error)
	GetCustomerPersonalByEmailAndNik(email, nik string) (data entity.CallbackData, err error)
	UpdateStatusDigisignActivation(email, nik, prospectID string, data []entity.TrxDetail) error
	SaveTrx(data []entity.TrxDetail) (err error)
	GetCustomerPersonalByEmail(documentID string) (data entity.CallbackData, err error)
	GetDigisignDummy(email string, action string) (data entity.DigisignDummy, err error)
	GetTrxMetadata(prospectID string) (data entity.TrxMetadata, err error)
	UpdateStatusDigisignSignDoc(data entity.TrxDetail) error
	SaveToWorker(data []entity.TrxWorker) (err error)
	GetDataWorker(prospectID string) (data entity.DataWorker, err error)
	SaveToTrxDigisign(data entity.TrxDigisign) (err error)
	GetTrxStatus(prospectID string) (status entity.TrxStatus, err error)
	GetLinkTrxDegisign(prospectID, action string) (data entity.TrxDigisign, err error)
	CheckWorker1209(prospectID string) (resultWorker int)
	CheckSND(prospectID string) (resultWorker int)
}
