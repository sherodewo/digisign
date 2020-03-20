package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type SignDocument struct {
	ID                   string             `gorm:"column:id;primary_key:true"`
	UserID               string             `gorm:"type:varchar(80);column:user_id"`
	DocumentID           string             `gorm:"type:varchar(255);column:document_id"`
	SignDocumentResult   SignDocumentResult `gorm:"forigenkey:SignDocumentResultID"`
	SignDocumentResultID string             `gorm:"column:sign_document_result_id"`
	EmailUser            string             `gorm:"type:varchar(100);column:email_user"`
	CreatedAt            time.Time          `gorm:"column:created_at"`
	UpdatedAt            time.Time          `gorm:"column:updated_at"`
}

func (c SignDocument) TableName() string {
	return "sign_documents"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SignDocument) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
