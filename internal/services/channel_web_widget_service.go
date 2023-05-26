package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type ChannelWebWidgetService struct{}

func NewChannelWebWidgetService() repository.ChannelWebWidgetRepository {
	return &ChannelWebWidgetService{}
}

func (*ChannelWebWidgetService) FindByWebsiteToken(websiteToken string) (*model.ChannelWebWidget, error) {
	var res model.ChannelWebWidget
	if err := config.GetDB().Orm().Debug().Model(&model.ChannelWebWidget{}).Where("website_token = ?", websiteToken).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
