package request

type AgentProfileRequest struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	EmployeeId string `json:"employeeId"`
	AccountId  string `json:"accountId"`
	Role       int32  `json:"role"`
	Manager    string `json:"manager"`
}
