package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type SignDocumentResult struct {
	ID             string       `gorm:"column:id;primary_key:true"`
	SignDocument   SignDocument `gorm:"forigenkey:SignDocumentID"`
	SignDocumentID string       `gorm:"column:sign_document_id"`
	Result         string       `gorm:"type:varchar(5);column:result"`
	Link           string       `gorm:"type:varchar(255);column:link"`
	CreatedAt      time.Time    `gorm:"column:created_at"`
	UpdatedAt      time.Time    `gorm:"column:updated_at"`
}

func (c SignDocumentResult) TableName() string {
	return "sign_document_result"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SignDocumentResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
