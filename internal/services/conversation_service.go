package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type ConversationService struct{}

func NewConversationService() repository.ConversationRepository {
	return &ConversationService{}
}

func (*ConversationService) FindByConversationId(conversationId string) (*model.Conversation, error) {
	var res model.Conversation
	if err := config.GetDB().Orm().Debug().Model(&model.Conversation{}).Where("uuid = ?", conversationId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}

func (*ConversationService) FindByAgentId(agentId string) (*model.Conversation, error) {
	var res model.Conversation
	if err := config.GetDB().Orm().Debug().Model(&model.Conversation{}).Where("agent_id = ?", agentId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
