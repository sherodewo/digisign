package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

type ActivationResult struct {
	ID           string     `gorm:"column:id;primary_key:true"`
	Activation   Activation `gorm:"forigenkey:ActivationID"`
	ActivationID string     `gorm:"column:activation_id"`
	Result       string     `gorm:"type:varchar(5);column:result"`
	Link         string     `gorm:"type:varchar(255);column:link"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (c ActivationResult) TableName() string {
	return "activation_result"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ActivationResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
