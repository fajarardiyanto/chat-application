package mapper

import (
	"github.com/fajarardiyanto/chat-application/internal/constant"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/model"
)

func AgentProfileMapper(agentProfile *model.AgentProfile, account *model.Account) (res *response.AgentProfileResponse) {
	res = &response.AgentProfileResponse{
		AgentId:    agentProfile.Uuid,
		FirstName:  agentProfile.FirstName,
		LastName:   agentProfile.LastName,
		Email:      agentProfile.Email,
		Phone:      agentProfile.Phone,
		Active:     agentProfile.Active,
		EmployeeId: agentProfile.EmployeeId,
		Account: response.AccountResponse{
			AccountId:   agentProfile.AccountId,
			AccountName: account.AccountName,
		},
		Role:       constant.AgentRole(agentProfile.Role),
		Supervisor: agentProfile.Manager,
		Photo:      agentProfile.Photo,
	}

	return res
}

func LoginMapper(agentProfile *model.AgentProfile, token string) (res *response.LoginResponse) {
	res = &response.LoginResponse{
		Token: token,
		AgentProfile: response.AgentProfileResponse{
			AgentId:    agentProfile.Uuid,
			FirstName:  agentProfile.FirstName,
			LastName:   agentProfile.LastName,
			Email:      agentProfile.Email,
			Phone:      agentProfile.Phone,
			Active:     agentProfile.Active,
			EmployeeId: agentProfile.EmployeeId,
			Account: response.AccountResponse{
				AccountId: agentProfile.AccountId,
			},
			Role:       constant.AgentRole(agentProfile.Role),
			Supervisor: agentProfile.Manager,
			Photo:      agentProfile.Photo,
		},
	}

	return res
}

func AllAgentProfileMapper(agentProfiles []model.AgentProfile, account *model.Account) (res []response.AgentProfileResponse) {
	for _, agentProfile := range agentProfiles {
		res = append(res, response.AgentProfileResponse{
			AgentId:    agentProfile.Uuid,
			FirstName:  agentProfile.FirstName,
			LastName:   agentProfile.LastName,
			Email:      agentProfile.Email,
			Phone:      agentProfile.Phone,
			Active:     agentProfile.Active,
			EmployeeId: agentProfile.EmployeeId,
			Account: response.AccountResponse{
				AccountId:   agentProfile.AccountId,
				AccountName: account.AccountName,
			},
			Role:       constant.AgentRole(agentProfile.Role),
			Supervisor: agentProfile.Manager,
			Photo:      agentProfile.Photo,
		})
	}

	return res
}
