package database

import (
	"fmt"
	"los-int-digisign/shared/config"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

func OpenLos() (*gorm.DB, error) {

	user, pwd, host, port, database := config.GetLosDB()

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		user, pwd, host, port, database,
	)

	db, err := gorm.Open("mssql", connString)
	if err != nil {
		return nil, err
	}

	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTION"))
	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTION"))

	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxOpen)
	db.DB().SetConnMaxLifetime(time.Hour)
	db.LogMode(config.IsDevelopment)

	return db, nil
}
