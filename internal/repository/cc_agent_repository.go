package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type CCAgentRepository interface {
	FindCCAgentByAgentId(agentId string) (*model.CCAgent, error)
}
