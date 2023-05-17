package repo

import (
	"github.com/fajarardiyanto/chat-application/internal/model"
)

type AgentProfileRepository interface {
	FindAgentProfileByEmail(email string) (*model.AgentProfile, error)
	Register(req model.AgentProfile) (*model.AgentProfile, error)
}
