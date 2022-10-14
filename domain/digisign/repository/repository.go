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
	WHERE tm.ProspectID = '%s' AND is_active = 1`, prospectID)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetCustomerPersonalByEmailAndNik(email, nik string) (data entity.CallbackData, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT TOP 1 cp.ProspectID, tm.redirect_success_url, tm.redirect_failed_url, ts.decision, DATEDIFF(minute, td.created_at, GETDATE()) AS diff_time
	 FROM customer_personal cp WITH (nolock) 
	 INNER JOIN trx_details td WITH (nolock) ON cp.ProspectID = td.ProspectID
	 INNER JOIN trx_metadata tm WITH (nolock) ON cp.ProspectID = tm.ProspectID
	 INNER JOIN trx_status ts WITH (nolock) ON cp.ProspectID = ts.ProspectID 
	 WHERE IDNumber = '%s' AND Email = '%s' AND td.source_decision = 'ACT' AND td.created_at >= DATEADD(minute, -10, GETDATE())
	 ORDER BY cp.created_at DESC`, nik, email)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetTrxMetadata(prospectID string) (data entity.TrxMetadata, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT TOP 1 redirect_url FROM trx_metadata WITH (nolock) WHERE ProspectID = '%s'`, prospectID)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetCustomerPersonalByEmail(documentID string) (data entity.CallbackData, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT TOP 1 cp.ProspectID, tm.redirect_success_url, tm.redirect_failed_url, ts.decision, DATEDIFF(minute, td.created_at, GETDATE()) AS diff_time
	FROM customer_personal cp WITH (nolock) 
	INNER JOIN trx_details td WITH (nolock) ON cp.ProspectID = td.ProspectID
	INNER JOIN trx_metadata tm WITH (nolock) ON cp.ProspectID = tm.ProspectID
	INNER JOIN trx_status ts WITH (nolock) ON cp.ProspectID = ts.ProspectID 
	WHERE CAST(info AS VARCHAR(30)) = '%s.pdf' AND td.source_decision = 'SID' AND td.created_at >= DATEADD(minute, -20, GETDATE())
	ORDER BY cp.created_at DESC`, documentID)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetAgreementNo(prospectID string) (data entity.TrxDetail, err error) {

	if err = r.db.Raw(fmt.Sprintf("SELECT info FROM trx_details WITH (nolock) WHERE ProspectID = '%s' AND source_decision = 'SND' AND decision = 'PAS'", prospectID)).Scan(&data).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) UpdateStatusDigisignActivation(email, nik, prospectID string, data []entity.TrxDetail) error {

	var digisignCustomer entity.DigisignCustomer

	result := r.db.Raw(fmt.Sprintf("SELECT ProspectID FROM digisign_customer WITH (nolock) WHERE Email = '%s' AND IDNumber = '%s' AND activation = 1", email, nik)).Scan(&digisignCustomer)

	if result.RowsAffected == 0 {

		return r.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Table("digisign_customer").Where("Email = ? AND IDNumber = ?", email, nik).Updates(&entity.DigisignCustomer{
				Activation:         1,
				DatetimeActivation: time.Now(),
			}).Error; err != nil {
				return err
			}

			if err := tx.Table("trx_details").Where("ProspectID = ? AND source_decision = ?", prospectID, "ACT").Updates(&entity.TrxDetail{
				SourceDecision: "ACT",
				Activity:       "PRCD",
				Decision:       "PAS",
				NextStep:       "SND",
			}).Error; err != nil {
				return err
			}

			latestDetails := data[len(data)-1]

			for _, details := range data {
				if err := tx.Create(&details).Error; err != nil {
					return err
				}
			}

			if err := tx.Table("trx_status").Where("ProspectID = ?", latestDetails.ProspectID).Updates(&entity.TrxStatus{
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

	} else {
		return nil
	}

}

func (r repoHandler) UpdateStatusDigisignSignDoc(data entity.TrxDetail, doc entity.TteDocPk) error {

	var details entity.TrxDetail

	result := r.db.Raw(fmt.Sprintf("SELECT ProspectID FROM trx_details WITH (nolock) WHERE ProspectID = '%s' AND rule_code = '%s'", data.ProspectID, data.RuleCode)).Scan(&details)

	if result.RowsAffected == 0 {
		return r.db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Table("trx_details").Where("ProspectID = ? AND source_decision = ?", data.ProspectID, "SID").Updates(&entity.TrxDetail{
				Activity: "PRCD", Decision: "PAS",
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&data).Error; err != nil {
				return err
			}

			if err := tx.Table("trx_status").Where("ProspectID = ? AND source_decision = ?", data.ProspectID, "SID").Updates(&entity.TrxStatus{
				Activity: data.Activity, Decision: data.Decision, RuleCode: data.RuleCode, NextStep: nil,
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&doc).Error; err != nil {
				return err
			}

			return nil
		})

	} else {
		return nil
	}

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

func (r repoHandler) SaveToWorker(data []entity.TrxWorker) (err error) {

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, worker := range data {
			if err := r.db.Create(&worker).Error; err != nil {
				return err
			}
		}

		return nil
	})

}

func (r repoHandler) GetDataWorker(prospectID string) (data entity.DataWorker, err error) {

	if err = r.db.Raw(fmt.Sprintf(`SELECT transaction_type, AF, tenor_limit, customer_id, callback_url FROM trx_master tm WITH(nolock)
	 INNER JOIN customer_kreditmu ck WITH (nolock) ON tm.ProspectID = ck.ProspectID
	 INNER JOIN trx_apk ta WITH (nolock) ON tm.ProspectID = ta.ProspectID
	 INNER JOIN trx_metadata tda WITH (nolock) ON tm.ProspectID = tda.ProspectID
	 WHERE tm.ProspectID = '%s'`, prospectID)).Scan(&data).Error; err != nil {
		return
	}
	return
}

func (r repoHandler) CheckWorker1209(prospectID string) (resultWorker int) {

	var check entity.CheckWorker

	result := r.db.Raw(fmt.Sprintf("SELECT ProspectID FROM trx_worker WITH (nolock) WHERE ProspectID = '%s' AND [action] = 'CALLBACK_STATUS_1209'", prospectID)).Scan(&check)

	resultWorker = int(result.RowsAffected)

	return
}

func (r repoHandler) CheckSND(prospectID string) (resultWorker int) {

	var check entity.CheckWorker

	result := r.db.Raw(fmt.Sprintf("SELECT ProspectID FROM trx_details WITH (nolock) WHERE ProspectID = '%s' AND source_decision = 'SND'", prospectID)).Scan(&check)

	resultWorker = int(result.RowsAffected)

	return
}

func (r repoHandler) SaveToTrxDigisign(data entity.TrxDigisign) (err error) {

	var check entity.TrxDigisign

	result := r.db.Raw(fmt.Sprintf("SELECT ProspectID FROM trx_digisign WITH (nolock) WHERE ProspectID = '%s' AND activity = '%s'", data.ProspectID, data.Activity)).Scan(&check)

	if result.RowsAffected == 0 {
		if err = r.db.Create(&data).Error; err != nil {
			return
		}

	}

	return
}

func (r repoHandler) GetTrxStatus(prospectID string) (status entity.TrxStatus, err error) {

	if err = r.db.Raw(fmt.Sprintf("SELECT * FROM trx_status WITH (nolock) WHERE ProspectID = '%s'", prospectID)).Scan(&status).Error; err != nil {
		return
	}

	return
}

func (r repoHandler) GetLinkTrxDegisign(prospectID, action string) (data entity.TrxDigisign, err error) {

	if err = r.db.Raw(fmt.Sprintf("SELECT TOP 1 link FROM trx_digisign WITH (nolock) WHERE ProspectID = '%s' AND activity = '%s' ORDER BY created_at DESC", prospectID, action)).Scan(&data).Error; err != nil {
		return
	}

	return
}
