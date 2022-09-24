package model

import "time"

type CustomerPersonal struct {
	ProspectID        string     `json:"prospect_id"`
	IDType            string     `json:"id_type"`
	IDNumber          string     `json:"id_number"`
	IDTypeIssueDate   *time.Time `json:"id_type_issue_date"`
	ExpiredDate       *time.Time `json:"expired_date"`
	LegalName         string     `json:"legal_name"`
	BirthPlace        string     `json:"birth_place"`
	BirthDate         string     `json:"birth_date"`
	SurgateMotherName string     `json:"surgate_mother_name"`
	Gender            string     `json:"gender"`
	PersonalNPWP      string     `json:"personal_npwp"`
	MobilePhone       string     `json:"mobile_phone"`
	Email             string     `json:"email"`
	HomeStatus        string     `json:"home_status"`
	StaySinceYear     string     `json:"stay_since_year"`
	StaySinceMonth    string     `json:"stay_since_month"`
	Education         string     `json:"education"`
	MaritalStatus     string     `json:"marital_status"`
	NumOfDependence   int        `json:"num_of_dependence"`
	LivingCostAmount  float64    `json:"living_cost_amount"`
	Religion          string     `json:"religion"`
}

func (c CustomerPersonal) TableName() string {
	return "customer_personal"
}
