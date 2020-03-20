package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type ActivationCallback struct {
	ID                 string           `gorm:"column:id;primary_key:true"`
	Email              string           `gorm:"type:varchar(100);column:email"`
	Result             string           `gorm:"type:varchar(5);column:result"`
	Notif              string           `gorm:"type:varchar(100);column:notif"`
	CreatedAt          time.Time        `gorm:"column:created_at"`
	UpdatedAt          time.Time        `gorm:"column:updated_at"`
}

func (c ActivationCallback) TableName() string {
	return "activation_callback"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ActivationCallback) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
