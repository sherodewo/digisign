package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type SignDocumentCallback struct {
	ID             string    `gorm:"column:id;primary_key:true"`
	DocumentID     string    `gorm:"type:varchar(255);column:document_id"`
	Email          string    `gorm:"type:varchar(100);column:email"`
	StatusDocument string    `gorm:"type:varchar(50);column:status_document"`
	Result         string    `gorm:"type:varchar(5);column:result"`
	Notif          string    `gorm:"type:varchar(50);column:notif"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (c SignDocumentCallback) TableName() string {
	return "sign_document_callback"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SignDocumentCallback) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
