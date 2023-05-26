package request

type RegisterContactRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Gender           string `json:"gender"`
	DateOfBirth      string `json:"date_of_birth"`
	MotherMaidenName string `json:"mother_maiden_name"`
	Note             string `json:"note"`
}
