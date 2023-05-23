package repository

import (
	"github.com/fajarardiyanto/chat-application/internal/model"
)

type AgentProfileRepository interface {
	FindAgentProfileByEmail(email string) (*model.AgentProfile, error)
	FindAgentProfileById(id string) (*model.AgentProfile, error)
	FindAgentProfileByAccountId(accountId string) ([]model.AgentProfile, error)
	UpdateAgentProfileById(req model.AgentProfile) error
	DeleteAgentProfileById(agentId string) error
	Register(req model.AgentProfile) (*model.AgentProfile, error)
}
