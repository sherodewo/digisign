package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"kpdigisign/model"
	"os"
)

var dataBase *gorm.DB

// New : err check
func NewDb() *gorm.DB {

	connection := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USERNAME") +
		" dbname=" + os.Getenv("DB_NAME") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" sslmode=" + os.Getenv("DB_SSL")
	db, err := gorm.Open(os.Getenv("DB_DRIVER"), connection)

	if err != nil {
		fmt.Println("Error DB: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	dataBase = db
	return db
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
