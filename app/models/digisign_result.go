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
	LosID        string `gorm:"column:los_request_id"`
	Result       string `gorm:"type:varchar(100);column:result"`
	Notif        string `gorm:"type:varchar(100);column:notif"`
	JsonResponse string `gorm:"type:varchar(255);column:json_response"`
	//Name            bool      `gorm:"column:name"`
	//BirthPlace      bool      `gorm:"column:birth_place"`
	//BirthDate       bool      `gorm:"column:birth_date"`
	//Address         string    `gorm:"type:varchar(100);column:address"`
	//Info            string    `gorm:"type:varchar(100);column:info"`
	//EmailRegistered string    `gorm:"type:varchar(100);column:email_registered"`
	//SelfieMatch     bool      `gorm:"type:varchar(100);column:selfie_match"`
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
