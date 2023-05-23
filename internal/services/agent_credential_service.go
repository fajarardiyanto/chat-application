package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
)

type AgentCredentialService struct{}

func NewAgentCredentialService() repository.AgentCredentialRepository {
	return &AgentCredentialService{}
}

func (*AgentCredentialService) FindAgentCredentialByUsername(email string) (*model.AgentCredential, error) {
	var res model.AgentCredential
	if err := config.GetDB().Orm().Debug().Model(&model.AgentCredential{}).Where("username = ?", email).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (*AgentCredentialService) SetPassword(req request.LoginRequest) error {
	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	salt := utils.SaltPassword(password)

	agentCredential := model.AgentCredential{
		Password: string(password),
		SaltKey:  salt,
		UserName: req.Email,
	}

	if err = config.GetDB().Orm().Debug().Model(&model.AgentCredential{}).Create(&agentCredential).Error; err != nil {
		return err
	}
	return nil
}
