package model

import "time"

type TrxDetail struct {
	ProspectID     string    `gorm:"type:varchar(50);column:ProspectID;primary_key:true"`
	StatusProcess  string    `gorm:"type:varchar(3);column:status_process"`
	Activity       string    `gorm:"type:varchar(4);column:activity"`
	Decision       string    `gorm:"type:varchar(3);column:decision"`
	RuleCode       string    `gorm:"type:varchar(4);column:rule_code"`
	SourceDecision string    `gorm:"type:varchar(3);column:source_decision"`
	NextStep       string    `gorm:"type:varchar(3);column:next_step"`
	Type           string    `gorm:"type:varchar(3);column:next_step"`
	Info           string    `gorm:"column:info"`
	CreatedBy      string    `gorm:"type:varchar(100);column:created_by"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (c *TrxDetail) TableName() string {
	return "trx_details"
}