package services

import (
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
)

type AccountService struct{}

func NewAccountService() repo.AccountRepository {
	return &AccountService{}
}

func (a *AccountService) FindAccountByAccountId(accountId string) (*model.Account, error) {
	var res model.Account
	if err := config.GetDB().Orm().Debug().Model(&model.Account{}).Where("uuid = ?", accountId).First(&res).Error; err != nil {
		return nil, err
	}

	return &res, nil
}
