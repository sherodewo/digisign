package send_document

import (
	"github.com/jinzhu/gorm"
	"kpdigisign/model"
)

type Repository interface {
	FindAll() ([]model.SendDocument, error)
	FindById(string) (model.SendDocument, error)
	Save(model.SendDocument) (model.SendDocument, error)
	Update(model.SendDocument) (model.SendDocument, error)
	Delete(model.SendDocument) error
}

type sendDocumentRepository struct {
	*gorm.DB
}

func NewSendDocumentRepository(db *gorm.DB) Repository {
	return &sendDocumentRepository{DB: db}
}

func (r sendDocumentRepository) FindAll() ([]model.SendDocument, error) {
	var entities []model.SendDocument
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r sendDocumentRepository) FindById(id string) (model.SendDocument, error) {
	var entity model.SendDocument
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r sendDocumentRepository) Save(entity model.SendDocument) (model.SendDocument, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r sendDocumentRepository) Update(entity model.SendDocument) (model.SendDocument, error) {
	err := r.DB.Model(model.SendDocument{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r sendDocumentRepository) Delete(entity model.SendDocument) error {
	return r.DB.Delete(&entity).Error
}
