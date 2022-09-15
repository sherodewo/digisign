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
