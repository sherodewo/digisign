package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

// DigisignRegistrationResult Struct
type DigisignRegistrationResult struct {
	ID        	string    	`gorm:"column:id;primary_key:true"`
	LosRequest LosRequest    	`gorm:"forigenkey:LosRequestID"`
	LosRequestID string    `gorm:"column:los_request_id"`
	Result  	string    	`gorm:"type:varchar(100);column:result"`
	Notif  		string    	`gorm:"type:varchar(100);column:notif"`
	Name  		bool    	`gorm:"column:name"`
	BirthPlace  bool    	`gorm:"column:birth_place"`
	BirthDate  	bool    	`gorm:"column:birth_date"`
	Address  	string    	`gorm:"type:varchar(100);column:address"`
	Info  		string    	`gorm:"type:varchar(100);column:info"`
	CreatedAt 	time.Time 	`gorm:"column:created_at"`
	UpdatedAt 	time.Time 	`gorm:"column:updated_at"`
}

func (c DigisignRegistrationResult) TableName() string {
	return "digisign_registration_result"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *DigisignRegistrationResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
