package response

type AgentProfileResponse struct {
	AgentId    string          `json:"agentId"`
	FirstName  string          `json:"firstName"`
	LastName   string          `json:"lastName"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	Active     bool            `json:"active"`
	EmployeeId string          `json:"employeeId"`
	Account    AccountResponse `json:"account"`
	Role       string          `json:"role"`
	Supervisor string          `json:"supervisor"`
	Photo      string          `json:"photo"`
}
