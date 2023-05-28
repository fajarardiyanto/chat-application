package repository

import "github.com/fajarardiyanto/chat-application/internal/model"

type ConversationRepository interface {
	FindByConversationId(conversationId string) (*model.Conversation, error)
	FindByAgentId(agentId string) (*model.Conversation, error)
	FindByContactInboxId(contactInboxId string) (*model.Conversation, error)
	CreateConversation(request model.Conversation) error
}
