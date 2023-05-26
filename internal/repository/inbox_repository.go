package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type InboxRepository interface {
	FindByChannelId(channelId string) (*model.Inbox, error)
}
