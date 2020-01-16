package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

// DigisignRegistrationResult Struct
type DigisignResult struct {
	ID           string `gorm:"column:id;primary_key:true"`
	Los          Los    `gorm:"forigenkey:LosID"`
	LosID        string `gorm:"column:los_id"`
	Result       string `gorm:"type:varchar(100);column:result"`
	Notif        string `gorm:"type:varchar(100);column:notif"`
	RefTrx       string `gorm:"type:varchar(100);column:ref_trx"`
	JsonResponse string `gorm:"type:varchar(100);column:json_response"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (c DigisignResult) TableName() string {
	return "digisign_result"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *DigisignResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
