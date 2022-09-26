package activation

import (
	"github.com/jinzhu/gorm"
	"los-int-digisign/model"
)

type Repository interface {
	FindAll() ([]model.Activation, error)
	FindById(string) (model.Activation, error)
	Save(model.Activation) (model.Activation, error)
	SaveCallback(callback model.ActivationCallback) (model.ActivationCallback, error)
	FindLastByCreatedAt(string, string) (model.Activation, error)
	Update(model.Activation) (model.Activation, error)
	Delete(model.Activation) error
	SaveDigisign(model.TrxDigisign) (model.TrxDigisign, error)
	FindCustomer(email string) (model.CustomerPersonal, error)
}

type activationRepository struct {
	DB *gorm.DB
	DBLos *gorm.DB
}

func NewActivationRepository(db *gorm.DB, dbLos *gorm.DB) Repository {
	return &activationRepository{DB: db, DBLos: dbLos}
}

func (r activationRepository) FindAll() ([]model.Activation, error) {
	var entities []model.Activation
	err := r.DB.Find(&entities).Error
	return entities, err
}

func (r activationRepository) FindById(id string) (model.Activation, error) {
	var entity model.Activation
	err := r.DB.Where("id = ?", id).First(&entity).Error
	return entity, err
}

func (r activationRepository) Save(entity model.Activation) (model.Activation, error) {
	err := r.DB.Create(&entity).Error
	return entity, err
}

func (r activationRepository) Update(entity model.Activation) (model.Activation, error) {
	err := r.DB.Model(model.Activation{ID: entity.ID}).UpdateColumns(&entity).Error
	return entity, err
}

func (r activationRepository) Delete(entity model.Activation) error {
	return r.DB.Delete(&entity).Error
}

func (r activationRepository) SaveCallback(callback model.ActivationCallback) (model.ActivationCallback, error) {
	err := r.DB.Create(&callback).Error
	return callback, err
}

func (r activationRepository) FindLastByCreatedAt(prospectId string, emailUser string) (model.Activation, error) {
	var entity model.Activation
	err := r.DB.Where("prospect_id =?", prospectId).Where("email_user =?", emailUser).Order("created_at DESC").First(&entity).Error
	return entity, err
}


func (r activationRepository) SaveDigisign(entity model.TrxDigisign) (model.TrxDigisign, error) {
	err := r.DBLos.Create(&entity).Error
	return entity, err
}

func (r activationRepository) FindCustomer(email string) (model.CustomerPersonal, error) {
	var entity model.CustomerPersonal
	err := r.DBLos.Raw("SELECT TOP 1 prospect_id WHERE email = ?", email).Error
	return entity, err
}