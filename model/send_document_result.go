package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

type SendDocumentResult struct {
	ID           string    `gorm:"column:id;primary_key:true"`
	Result       string    `gorm:"type:varchar(100);column:result"`
	Notif        string    `gorm:"type:varchar(100);column:notif"`
	RefTrx       string    `gorm:"type:varchar(100);column:ref_trx"`
	JsonResponse string    `gorm:"type:text;column:json_response"`
	JsonRequest  string    `gorm:"type:text;column:json_request"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (c SendDocumentResult) TableName() string {
	return "send_document_results"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SendDocumentResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
