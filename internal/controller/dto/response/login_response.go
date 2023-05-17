package response

type LoginResponse struct {
	Token        string               `json:"token"`
	AgentProfile AgentProfileResponse `json:"agentProfile"`
}
