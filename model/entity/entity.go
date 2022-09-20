package entity

import "time"

type SendDocData struct {
	BranchID  string `gorm:"column:BranchID"`
	LegalName string `gorm:"column:LegalName"`
	Email     string `gorm:"column:Email"`
	NameBM    string `gorm:"column:name"`
	EmailBM   string `gorm:"column:email_bm"`
	Kuser     string `gorm:"column:kuser"`
}

type CustomerPersonal struct {
	ProspectID        string      `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	IDType            string      `gorm:"type:varchar(30);column:IDType"`
	IDNumber          string      `gorm:"type:varchar(100);column:IDNumber"`
	IDTypeIssueDate   interface{} `gorm:"column:IDTypeIssuedDate"`
	ExpiredDate       interface{} `gorm:"column:ExpiredDate"`
	LegalName         string      `gorm:"type:varchar(100);column:LegalName"`
	FullName          string      `gorm:"type:varchar(100);column:FullName"`
	BirthPlace        string      `gorm:"type:varchar(100);column:BirthPlace"`
	BirthDate         time.Time   `gorm:"column:BirthDate"`
	SurgateMotherName string      `gorm:"type:varchar(100);column:SurgateMotherName"`
	Gender            string      `gorm:"type:varchar(10);column:Gender"`
	PersonalNPWP      string      `gorm:"type:varchar(255);column:PersonalNPWP"`
	MobilePhone       string      `gorm:"type:varchar(14);column:MobilePhone"`
	Email             string      `gorm:"type:varchar(100);column:Email"`
	HomeStatus        string      `gorm:"type:varchar(20);column:HomeStatus"`
	StaySinceYear     string      `gorm:"type:varchar(10);column:StaySinceYear"`
	StaySinceMonth    string      `gorm:"type:varchar(10);column:StaySinceMonth"`
	Education         string      `gorm:"type:varchar(50);column:Education"`
	MaritalStatus     string      `gorm:"type:varchar(10);column:MaritalStatus"`
	NumOfDependence   int         `gorm:"column:NumOfDependence"`
	LivingCostAmount  float64     `gorm:"column:LivingCostAmount"`
	Religion          string      `gorm:"type:varchar(30);column:Religion"`
	CreatedAt         time.Time   `gorm:"column:created_at"`
}

func (c *CustomerPersonal) TableName() string {
	return "customer_personal"
}

type TrxDetail struct {
	ProspectID     string      `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	StatusProcess  string      `gorm:"type:varchar(3);column:status_process"`
	Activity       string      `gorm:"type:varchar(4);column:activity"`
	Decision       string      `gorm:"type:varchar(3);column:decision"`
	RuleCode       interface{} `gorm:"type:varchar(4);column:rule_code"`
	SourceDecision string      `gorm:"type:varchar(3);column:source_decision"`
	NextStep       interface{} `gorm:"type:varchar(3);column:next_step"`
	Type           interface{} `gorm:"type:varchar(3);column:type"`
	Info           interface{} `gorm:"type:text;column:info"`
	CreatedBy      string      `gorm:"type:varchar(100);column:created_by"`
	CreatedAt      time.Time   `gorm:"column:created_at"`
}

func (c *TrxDetail) TableName() string {
	return "trx_details"
}

type TrxStatus struct {
	ProspectID     string      `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	StatusProcess  string      `gorm:"type:varchar(3);column:status_process"`
	Activity       string      `gorm:"type:varchar(4);column:activity"`
	Decision       string      `gorm:"type:varchar(3);column:decision"`
	RuleCode       interface{} `gorm:"type:varchar(4);column:rule_code"`
	SourceDecision string      `gorm:"type:varchar(3);column:source_decision"`
	NextStep       interface{} `gorm:"type:varchar(3);column:next_step"`
	CreatedAt      time.Time   `gorm:"column:created_at"`
	Reason         interface{} `gorm:"type:varchar(255);column:reason"`
}

func (c *TrxStatus) TableName() string {
	return "trx_status"
}

type DigisignDummy struct {
	Email    string `gorm:"column:email"`
	Response string `gorm:"column:response"`
}

func (c *DigisignDummy) TableName() string {
	return "digisign_dummy"
}

type DigisignCustomer struct {
	IDNumber           string      `gorm:"type:varchar(40);column:IDNumber"`
	LegalName          string      `gorm:"type:varchar(100);column:LegalName"`
	BirthDate          string      `gorm:"column:BirthDate"`
	SurgateMotherName  string      `gorm:"type:varchar(100);column:SurgateMotherName"`
	MobilePhone        string      `gorm:"type:varchar(20);column:MobilePhone"`
	Email              string      `gorm:"type:varchar(100);column:Email"`
	Register           int         `gorm:"column:register"`
	DatetimeRegister   interface{} `gorm:"column:datetime_register"`
	Activation         int         `gorm:"column:activation"`
	DatetimeActivation interface{} `gorm:"column:datetime_activation"`
	ProspectID         string      `gorm:"type:varchar(20);column:ProspectID"`
}

func (c *DigisignCustomer) TableName() string {
	return "digisign_customer"
}
