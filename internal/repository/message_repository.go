package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type MessageRepository interface {
	CreateMessage(message model.Message) error
}
