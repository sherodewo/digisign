package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type RegistrationResult struct {
	ID           string    `gorm:"column:id;primary_key:true"`
	Result       string    `gorm:"type:varchar(5);column:result"`
	Notif        string    `gorm:"type:varchar(150);column:notif"`
	RefTrx       string    `gorm:"type:varchar(200);column:ref_trx"`
	KodeUser     string    `gorm:"type:varchar(255);column:kode_user"`
	JsonResponse string    `gorm:"type:text;column:json_response"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (m *RegistrationResult) TableName() string {
	return "registration_results"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (m *RegistrationResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
