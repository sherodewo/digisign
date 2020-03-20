package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type SendDocument struct {
	ID                   string             `gorm:"column:id;primary_key:true"`
	UserID               string             `gorm:"type:varchar(80);column:user_id"`
	DocumentID           string             `gorm:"type:varchar(20);column:document_id"`
	SendDocumentResult   SendDocumentResult `gorm:"forigenkey:SendDocumentResultID"`
	SendDocumentResultID string             `gorm:"column:send_document_result_id"`
	Payment              string             `gorm:"type:varchar(1);column:payment"`
	SendTo               string             `gorm:"column:send_to"`
	ReqSign              string             `gorm:"column:req_sign"`
	CreatedAt            time.Time          `gorm:"column:created_at"`
	UpdatedAt            time.Time          `gorm:"column:updated_at"`
}

func (c *SendDocument) TableName() string {
	return "send_documents"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *SendDocument) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
