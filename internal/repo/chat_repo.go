package repo

import (
	"github.com/fajarardiyanto/chat-application/internal/model"
	"net/http"
)

type ChatRepository interface {
	CreateChat(req model.Chat) (*model.Chat, error)
	GetChat(from, to string) ([]model.Chat, error)
	SaveFileChat(r *http.Request) (*model.FileModel, error)
}
