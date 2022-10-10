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

type CallbackData struct {
	ProspectID         string  `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	Decision           string  `gorm:"type:varchar(3);column:decision"`
	RedirectSuccessUrl string  `gorm:"type:varchar(250);column:redirect_success_url"`
	RedirectFailedUrl  string  `gorm:"type:varchar(250);column:redirect_failed_url"`
	DiffTime           float64 `gorm:"column:diff_time"`
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

type TrxMetadata struct {
	ProspectID         string    `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	CustomerIp         string    `gorm:"type:varchar(15);column:customer_ip"`
	CustomerLat        string    `gorm:"type:varchar(10);column:customer_lat"`
	CustomerLong       string    `gorm:"type:varchar(10);column:customer_long"`
	CallbackUrl        string    `gorm:"type:varchar(250);column:callback_url"`
	RedirectSuccessUrl string    `gorm:"type:varchar(250);column:redirect_success_url"`
	RedirectFailedUrl  string    `gorm:"type:varchar(250);column:redirect_failed_url"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}

func (c *TrxMetadata) TableName() string {
	return "trx_metadata"
}

type TrxWorker struct {
	ProspectID      string      `gorm:"type:varchar(20);column:ProspectID;primary_key:true"`
	Activity        string      `gorm:"type:varchar(10);column:activity"`
	EndPointTarget  string      `gorm:"type:varchar(100);column:endpoint_target"`
	EndPointMethod  string      `gorm:"type:varchar(10);column:endpoint_method"`
	Payload         string      `gorm:"type:text;column:payload"`
	Header          string      `gorm:"type:text;column:header"`
	ResponseTimeout int         `gorm:"column:response_timeout"`
	APIType         string      `gorm:"type:varchar(3);column:api_type"`
	MaxRetry        int         `gorm:"max_retry"`
	CountRetry      int         `gorm:"count_retry"`
	CreatedAt       time.Time   `gorm:"column:created_at"`
	Category        string      `gorm:"type:varchar(30);column:category"`
	Action          string      `gorm:"type:varchar(50);column:action"`
	StatusCode      string      `gorm:"type:varchar(4);column:status_code"`
	Sequence        interface{} `gorm:"column:sequence"`
}

func (c *TrxWorker) TableName() string {
	return "trx_worker"
}

type DataWorker struct {
	TransactionType string  `gorm:"type:varchar(30);column:transaction_type"`
	CustomerID      int     `gorm:"column:customer_id"`
	AF              float64 `gorm:"column:AF"`
	TenorLimit      int     `gorm:"column:tenor_limit"`
	CallbackUrl     string  `gorm:"column:callback_url"`
}

type CheckWorker struct {
	ProspectID string `gorm:"column:ProspectID"`
	Action     string `gorm:"column:action"`
}

type TrxDigisign struct {
	ProspectID string      `gorm:"type:varchar(20);column:ProspectID"`
	Response   string      `gorm:"type:varchar(2000);column:response"`
	Activity   string      `gorm:"type:varchar(20);column:activity"`
	Link       interface{} `gorm:"type:varchar(2000);column:link"`
	CreatedAt  time.Time   `gorm:"column:created_at"`
}

func (c *TrxDigisign) TableName() string {
	return "trx_digisign"
}

type TteDocPk struct {
	ID          string    `gorm:"type:varchar(50);column:id"`
	ProspectID  string    `gorm:"type:varchar(20);column:prospect_id"`
	NoAgreement string    `gorn:"type:varchar(50);column:no_agreement"`
	DocPKUrl    string    `gorm:"type:varchar(255);column:doc_pk_url"`
	Tipe        string    `gorm:"type:varchar(10);column:tipe"`
	FilePath    string    `gorm:"type:varchar(255);column:file_path"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (c *TteDocPk) TableName() string {
	return "tte_doc_pk"
}
