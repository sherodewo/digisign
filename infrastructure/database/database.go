package database

import (
	"fmt"
	"log"
	"los-int-digisign/model"
	"los-int-digisign/utils"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dataBase *gorm.DB

// New : err check
func NewDb() (*gorm.DB, error) {

	decryptDbHost, err := utils.DecryptCredential(os.Getenv("DB_HOST"))

	decryptDbPort, err := utils.DecryptCredential(os.Getenv("DB_PORT"))

	decryptDbPassword, err := utils.DecryptCredential(os.Getenv("DB_PASSWORD"))

	decryptDbUsername, err := utils.DecryptCredential(os.Getenv("DB_USERNAME"))

	decryptDbName, err := utils.DecryptCredential(os.Getenv("DB_NAME"))

	connection := "host=" + decryptDbHost +
		" port=" + decryptDbPort +
		" user=" + decryptDbUsername +
		" dbname=" + decryptDbName +
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


func NewLosDb() (*gorm.DB, error) {

	decryptDbHost, err := utils.DecryptCredential(os.Getenv("LOS_DB_HOST"))

	decryptDbPort, err := utils.DecryptCredential(os.Getenv("LOS_DB_PORT"))

	decryptDbPassword, err := utils.DecryptCredential(os.Getenv("LOS_DB_PASSWORD"))

	decryptDbUsername, err := utils.DecryptCredential(os.Getenv("LOS_DB_USERNAME"))

	decryptDbName, err := utils.DecryptCredential(os.Getenv("LOS_DB_NAME"))

	port, _ := strconv.Atoi(decryptDbPort)
	args := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s",
		decryptDbUsername,
		decryptDbPassword,
		decryptDbHost,
		port,
		decryptDbName,
	)
	db, err := gorm.Open("mssql", args)
	if err != nil {
		log.Println("Error", err)
	}
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)

	return db, err
}
