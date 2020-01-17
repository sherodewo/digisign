package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type Activation struct {
	ID        string    `gorm:"column:id;primary_key:true"`
	UserID    string    `gorm:"type:varchar(80);column:user_id"`
	EmailUser string    `gorm:"type:varchar(100);column:email_user"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (c Activation) TableName() string {
	return "activations"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *Activation) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
