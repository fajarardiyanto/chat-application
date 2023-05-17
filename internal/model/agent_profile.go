package model

type AgentProfile struct {
	Uuid       string `json:"agent_id" gorm:"column:uuid"`
	FirstName  string `json:"firstName" gorm:"column:first_name"`
	LastName   string `json:"lastName" gorm:"column:last_name"`
	Email      string `json:"email" gorm:"email"`
	Phone      string `json:"phone" gorm:"column:phone"`
	Active     bool   `json:"active" gorm:"column:active;default:true"`
	EmployeeId string `json:"employeeId" gorm:"column:employee_id"`
	AccountId  string `json:"account" gorm:"column:account_uuid"`
	Role       int32  `json:"role" gorm:"column:role"`
	Manager    string `json:"manager" gorm:"column:manager"`
	Photo      string `json:"photo" gorm:"column:photo"`
	Deleted    string `json:"deleted" gorm:"column:deleted;default:false"`
	Audit
}

func (*AgentProfile) TableName() string {
	return "agent_profiles"
}
