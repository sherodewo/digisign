package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
)

type ActivationResult struct {
	ID           string    `gorm:"column:id;primary_key:true"`
	Result       string    `gorm:"type:varchar(5);column:result"`
	Link         string    `gorm:"type:text;column:link"`
	JsonResponse string    `gorm:"type:text;column:json_response"`
	JsonRequest  string    `gorm:"type:text;column:json_request"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (c ActivationResult) TableName() string {
	return "activation_results"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *ActivationResult) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
