package sign_document

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"los-int-digisign/model"
)

type Repository interface {
	FindAll() ([]model.SignDocument, error)
	FindById(string) (model.SignDocument, error)
	Save(model.SignDocument) (model.SignDocument, error)
	SaveCallback(callback model.SignDocumentCallback) (model.SignDocumentCallback, error)
	Update(model.SignDocument) (model.SignDocument, error)
	Delete(model.SignDocument) error
	SaveDigisign(model.TrxDigisign) (model.TrxDigisign, error)
	FindCustomer(string) (model.TrxDetail, error)
}

type signDocumentRepository struct {
	*gorm.DB
}

func NewSignDocumentRepository(db *gorm.DB) Repository {
	return &signDocumentRepository{DB: db}
}

func (r signDocumentRepository) FindAll() ([]model.SignDocument, error) {
	var entities []model.SignDocument
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r signDocumentRepository) FindById(id string) (model.SignDocument, error) {
	var entity model.SignDocument
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r signDocumentRepository) Save(entity model.SignDocument) (model.SignDocument, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r signDocumentRepository) Update(entity model.SignDocument) (model.SignDocument, error) {
	err := r.DB.Model(model.SignDocument{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r signDocumentRepository) Delete(entity model.SignDocument) error {
	return r.DB.Delete(&entity).Error
}

func (r signDocumentRepository) SaveCallback(callback model.SignDocumentCallback) (model.SignDocumentCallback, error) {
	err := r.DB.Create(&callback).Error
	return callback, err
}

func (r signDocumentRepository) SaveDigisign(entity model.TrxDigisign) (model.TrxDigisign, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r signDocumentRepository) FindCustomer(docID string) (model.TrxDetail, error) {
	var entity model.TrxDetail
	err := r.DB.Raw(fmt.Sprintf(`
	SELECT
	* 
	FROM
	trx_details 
	WHERE
	status_process = 'ONP' 
	AND activity = 'PRGS' 
	AND rule_code IN ( '1206', '1207' ) 
	AND JSON_VALUE ( CAST ( info AS NVARCHAR ( MAX ) ), '$.document_id' ) = '%s'
	AND  created_at >= DATEADD(day, -2, GETDATE())`, docID)).Error
	return entity, err
}