package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/google/uuid"
	"time"
)

type ChatService struct{}

func NewChatService() repository.ChatRepository {
	return &ChatService{}
}

func (s *ChatService) CreateChat(req model.Chat) (*model.Chat, error) {
	req.ID = uuid.NewString()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := config.GetDB().Orm().Debug().Model(&model.Chat{}).Create(&req).Error; err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	return &req, nil
}

func (s *ChatService) GetChat(from, to string) ([]model.Chat, error) {
	var res []model.Chat
	if err := config.GetDB().Orm().Debug().Model(&model.Chat{}).Where("from_user = ? AND to_user = ? OR from_user = ? AND to_user = ?", from, to, to, from).Order("created_at desc").Find(&res).Error; err != nil {
		config.GetLogger().Error(err)
		return nil, err
	}

	return res, nil
}
