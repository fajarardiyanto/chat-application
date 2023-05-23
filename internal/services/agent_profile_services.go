package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"time"
)

type AgentProfileService struct{}

func NewAgentProfileService() repository.AgentProfileRepository {
	return &AgentProfileService{}
}

func (*AgentProfileService) FindAgentProfileByEmail(email string) (*model.AgentProfile, error) {
	var res model.AgentProfile
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Where("email = ? AND deleted = false", email).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (*AgentProfileService) FindAgentProfileById(id string) (*model.AgentProfile, error) {
	var res model.AgentProfile
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Where("uuid = ?", id).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (*AgentProfileService) FindAgentProfileByAccountId(accountId string) ([]model.AgentProfile, error) {
	var res []model.AgentProfile
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Where("account_uuid = ?", accountId).Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (*AgentProfileService) Register(req model.AgentProfile) (*model.AgentProfile, error) {
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Create(&req).Error; err != nil {
		return nil, err
	}

	return &req, nil
}

func (*AgentProfileService) UpdateAgentProfileById(req model.AgentProfile) error {
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Where("uuid = ?", req.Uuid).Updates(map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"phone":      req.Phone,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (*AgentProfileService) DeleteAgentProfileById(agentId string) error {
	data := model.AgentProfile{
		Deleted: true,
	}
	if err := config.GetDB().Orm().Debug().Model(&model.AgentProfile{}).Where("uuid = ?", agentId).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
