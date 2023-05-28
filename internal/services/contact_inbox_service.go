package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type ContactInboxService struct{}

func NewContactInboxService() repository.ContactInboxRepository {
	return &ContactInboxService{}
}

func (*ContactInboxService) CreateContactInbox(contactInbox model.ContactInbox) error {
	if err := config.GetDB().Orm().Debug().Model(&model.ContactInbox{}).Create(&contactInbox).Error; err != nil {
		return err
	}

	return nil
}

func (*ContactInboxService) FindBySourceId(sourceId string) (*model.ContactInbox, error) {
	var res model.ContactInbox
	if err := config.GetDB().Orm().Debug().Model(&model.ContactInbox{}).Where("source_id = ?", sourceId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
