package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
)

type ContactService struct{}

func NewContactService() repository.ContactRepository {
	return &ContactService{}
}

func (*ContactService) RegisterContact(contact model.Contact) error {
	if err := config.GetDB().Orm().Debug().Model(&model.Contact{}).Create(&contact).Error; err != nil {
		return err
	}

	return nil
}

func (*ContactService) FindByEmailAndAccountUuid(email string, accountId string) (*model.Contact, error) {
	var res model.Contact
	if err := config.GetDB().Orm().Debug().Model(&model.Contact{}).Where("email = ? AND account_id = ?", email, accountId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
