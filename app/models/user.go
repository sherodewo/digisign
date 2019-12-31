package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"time"
)

// User Struct
type User struct {
	ID        string    `gorm:"column:id;primary_key:true"`
	Username  string    `gorm:"type:varchar(100);column:username"`
	Email     string    `gorm:"type:varchar(50);column:email;unique"`
	Password  string    `gorm:"type:varchar(255);column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (c User) TableName() string {
	return "users"
}

// BeforeCreate - Lifecycle callback - Generate UUID before persisting
func (c *User) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", uuid.New().String()); err != nil {
		log.Fatal("Error UUID Generate")
	}
	return nil
}
