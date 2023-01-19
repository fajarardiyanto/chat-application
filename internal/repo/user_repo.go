package repo

import "github.com/fajarardiyanto/chat-application/internal/model"

type UserRepository interface {
	UserExist(username string) error
	Register(req model.UserModel) (*model.UserModel, error)
	GetUser() ([]model.UserModel, error)
	UpdateStatus(id string, status bool) error
	CheckUserLife(id string) bool
}
