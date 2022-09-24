package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// TrxDigisign Struct
type TrxDigisign struct {
	ProspectID string    `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	Response   string    `gorm:"type:varchar(2000);column:response"`
	Activity   string    `gorm:"type:varchar(20);column:activity"`
	Link       string    `gorm:"type:varchar(2000);column:link"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (c TrxDigisign) TableName() string {
	return "trx_digisign"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *TrxDigisign) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
