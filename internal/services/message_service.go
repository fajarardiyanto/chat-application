package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type MessageService struct{}

func NewMessageService() repository.MessageRepository {
	return &MessageService{}
}

func (m MessageService) CreateMessage(message model.Message) error {
	if err := config.GetDB().Orm().Debug().Model(&model.Message{}).Create(&message).Error; err != nil {
		return err
	}

	return nil
}
