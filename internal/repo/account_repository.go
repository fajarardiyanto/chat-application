package repo

import "github.com/fajarardiyanto/chat-application/internal/model"

type AccountRepository interface {
	FindAccountByAccountId(accountId string) (*model.Account, error)
}
