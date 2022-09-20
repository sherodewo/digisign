package repository

import (
	"fmt"
	"los-int-digisign/domain/digisign/interfaces"
	"los-int-digisign/model/entity"
	"time"

	"github.com/jinzhu/gorm"
)

type repoHandler struct {
	digisign *gorm.DB
	db       *gorm.DB
	dbLog    *gorm.DB
}

func NewRepository(digisign *gorm.DB, db *gorm.DB, dbLog *gorm.DB) interfaces.Repository {
	return &repoHandler{
		digisign: digisign,
		db:       db,
		dbLog:    dbLog,
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

	if err = r.db.Raw(fmt.Sprintf(`SELECT TOP 1 ProspectID FROM customer_personal WHERE IDNumber = '%s' AND Email = '%s'`, nik, email)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetCustomerPersonalByEmail(documentID string) (data entity.TrxDetail, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT TOP 1 ProspectID FROM trx_details WHERE CAST(info AS VARCHAR(30)) = '%s.pdf' AND source_decision = 'SID' AND created_at > DATEADD(day, -2, GETDATE()) `, documentID)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) UpdateStatusDigisignActivation(prospectID string) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("digisign_customer").Where("ProspectID = ?", prospectID).Update(&entity.DigisignCustomer{
			Activation:         1,
			DatetimeActivation: time.Now(),
		}).Error; err != nil {
			return err
		}

		if err := tx.Table("trx_details").Where("ProspectID = ? AND source_decision = ?", prospectID, "ACT").Update(&entity.TrxDetail{
			SourceDecision: "ACT",
			Activity:       "PRCD",
			Decision:       "PAS",
			NextStep:       "SND",
		}).Error; err != nil {
			return err
		}

		if err := tx.Table("trx_status").Where("ProspectID = ? AND source_decision = ?", prospectID, "ACT").Update(&entity.TrxStatus{
			Activity: "PRCD",
			Decision: "PAS",
			NextStep: "SND",
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r repoHandler) SaveTrx(data []entity.TrxDetail) (err error) {

	latestDetails := data[len(data)-1]

	return r.db.Transaction(func(tx *gorm.DB) error {

		for _, details := range data {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}

		if err := tx.Table("trx_status").Where("ProspectID = ?", latestDetails.ProspectID).Update(&entity.TrxStatus{
			ProspectID:     latestDetails.ProspectID,
			StatusProcess:  latestDetails.StatusProcess,
			Activity:       latestDetails.Activity,
			Decision:       latestDetails.Decision,
			RuleCode:       latestDetails.RuleCode,
			SourceDecision: latestDetails.SourceDecision,
			NextStep:       latestDetails.NextStep,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r repoHandler) GetDigisignDummy(email string, action string) (data entity.DigisignDummy, err error) {

	if err = r.dbLog.Raw(fmt.Sprintf("SELECT response FROM digisign_dummy WHERE email = '%s' AND action = '%s'", email, action)).Scan(&data).Error; err != nil {
		return
	}

	return
}
