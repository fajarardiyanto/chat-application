package repository

import (
	"github.com/fajarardiyanto/chat-application/internal/model"
)

type ChatRepository interface {
	CreateChat(req model.Chat) (*model.Chat, error)
	GetChat(from, to string) ([]model.Chat, error)
}
