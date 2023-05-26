package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type ChannelWebWidgetRepository interface {
	FindByWebsiteToken(websiteToken string) (*model.ChannelWebWidget, error)
}
