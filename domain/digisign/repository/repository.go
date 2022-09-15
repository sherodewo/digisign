package repository

import (
	"fmt"
	"los-int-digisign/domain/digisign/interfaces"
	"los-int-digisign/model/entity"

	"github.com/jinzhu/gorm"
)

type repoHandler struct {
	digisign *gorm.DB
	db       *gorm.DB
}

func NewRepository(digisign *gorm.DB, db *gorm.DB) interfaces.Repository {
	return &repoHandler{
		digisign: digisign,
		db:       db,
	}
}

func (r repoHandler) GetSendDocData(prospectID string) (data entity.SendDocData, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT tm.BranchID, cp.LegalName, cp.Email, ub.name, ub.email AS email_bm, ub.kuser
	FROM trx_master tm INNER JOIN customer_personal cp
	ON tm.ProspectID = cp.ProspectID
	LEFT JOIN user_bm ub
	ON tm.BranchID = ub.branch_id
	WHERE tm.ProspectID = '%s'`, prospectID)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetCustomerPersonalByEmailAndNik(email, nik string) (data entity.CustomerPersonal, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT ProspectID FROM customer_personal WHERE IDNumber = '%s' AND Email = '%s'`, nik, email)).Scan(&data).Error; err != nil {
		return
	}

	return
}
