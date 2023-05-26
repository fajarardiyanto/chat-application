package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type InboxService struct{}

func NewInboxService() repository.InboxRepository {
	return &InboxService{}
}

func (*InboxService) FindByChannelId(channelId string) (*model.Inbox, error) {
	var res model.Inbox
	if err := config.GetDB().Orm().Debug().Model(&model.Inbox{}).Where("channel_id = ?", channelId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
