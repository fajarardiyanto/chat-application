package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
)

type AgentProfileService struct{}

func NewAgentProfileService() repo.AgentProfileRepository {
	return &AgentProfileService{}
}

func (*AgentProfileService) FindAgentProfileByEmail(email string) (*model.AgentProfile, error) {
	var res model.AgentProfile
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Where("email = ?", email).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (*AgentProfileService) Register(req model.AgentProfile) (*model.AgentProfile, error) {
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Create(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}
