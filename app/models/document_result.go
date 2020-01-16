package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type DocumentResult struct {
	ID           string    `gorm:"column:id;primary_key:true"`
	Document     Document  `gorm:"forigenkey:DocumentID"`
	DocumentID   string    `gorm:"column:document_id"`
	Result       string    `gorm:"type:varchar(100);column:result"`
	Notif        string    `gorm:"type:varchar(100);column:notif"`
	JsonResponse string    `gorm:"type:varchar(255);column:json_response"`
	RefTrx       string    `gorm:"type:varchar(100);column:ref_trx"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (c DocumentResult) TableName() string {
	return "document_result"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *DocumentResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
