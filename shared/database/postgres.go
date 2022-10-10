package database

import (
	"fmt"
	"los-int-digisign/shared/config"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func OpenDigisign() (*gorm.DB, error) {

	user, pwd, host, port, database := config.DigisignDBCredential()

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		host, port, user, database, pwd,
	)

	db, err := gorm.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	maxIdle, _ := strconv.Atoi(os.Getenv("DIGISIGN_DB_MAX_IDLE_CONNECTION"))
	maxOpen, _ := strconv.Atoi(os.Getenv("DIGISIGN_DB_MAX_OPEN_CONNECTION"))

	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	db.LogMode(config.IsDevelopment)

	return db, nil
}
