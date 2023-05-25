package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type CCAgentService struct{}

func NewCCAgentService() repository.CCAgentRepository {
	return &CCAgentService{}
}

func (*CCAgentService) FindCCAgentByAgentId(agentId string) (*model.CCAgent, error) {
	var res model.CCAgent
	if err := config.GetDB().Orm().Debug().Model(&model.CCAgent{}).Where("agent_id = ?", agentId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
