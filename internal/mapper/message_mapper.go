package mapper

import (
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/model"
)

func MessageMapper(message model.Message) response.MessageResponse {
	return response.MessageResponse{
		Content:        message.Content,
		ConversationId: message.ConversationId,
	}
}
