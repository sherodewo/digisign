package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"kpdigisign/model"
	"kpdigisign/utils"
	"os"
)

var dataBase *gorm.DB

// New : err check
func NewDb() (*gorm.DB, error) {

	decryptDbHost, err := utils.DecryptCredential(os.Getenv("DB_HOST"))

	decryptDbPort, err := utils.DecryptCredential(os.Getenv("DB_PORT"))

	decryptDbPassword, err := utils.DecryptCredential(os.Getenv("DB_PASSWORD"))

	decryptDbUsername, err := utils.DecryptCredential(os.Getenv("DB_USERNAME"))


	connection := "host=" + decryptDbHost +
		" port=" + decryptDbPort +
		" user=" + decryptDbUsername +
		" dbname=" + os.Getenv("DB_NAME") +
		" password=" + decryptDbPassword +
		" sslmode=" + os.Getenv("DB_SSL")
	db, err := gorm.Open(os.Getenv("DB_DRIVER"), connection)

	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)

	return db, err
}

// GetLinkDb :
func GetLinkDb() *gorm.DB {
	return dataBase
}

// AutoMigrate : err check
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.RegistrationResult{},
		&model.Registration{},
		&model.SendDocumentResult{},
		&model.SendDocument{},
		&model.ActivationResult{},
		&model.Activation{},
		&model.SignDocumentResult{},
		&model.SignDocument{},
		&model.ActivationCallback{},
		&model.SignDocumentCallback{},
	)
}
