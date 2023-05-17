package repo

import (
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/model"
)

type AgentCredentialRepository interface {
	FindAgentCredentialByUsername(email string) (*model.AgentCredential, error)
	SetPassword(req request.LoginRequest) error
}
