package model

import (
	"time"
)

// TrxDigisign Struct
type TrxDigisign struct {
	ProspectID string    `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	Response   string    `gorm:"type:varchar(2000);column:response"`
	Activity   string    `gorm:"type:varchar(20);column:activity"`
	Link       string    `gorm:"type:varchar(2000);column:link"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (c TrxDigisign) TableName() string {
	return "trx_digisign"
}