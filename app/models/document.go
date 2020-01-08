package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type Document struct {
	ID         string    `gorm:"column:id;primary_key:true"`
	UserID     string    `gorm:"type:varchar(80);column:user_id"`
	DocumentID string    `gorm:"type:varchar(20);column:document_id"`
	Payment    string    `gorm:"type:varchar(1);column:password"`
	SendTo     string    `gorm:"type:varchar(150);column:send_to"`
	ReqSign    string    `gorm:"type:varchar(150);column:req_sign"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (c Document) TableName() string {
	return "documents"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Document) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
