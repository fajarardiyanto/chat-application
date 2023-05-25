package controller

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/commons"
	"github.com/fajarardiyanto/chat-application/internal/constant"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/exception"
	"github.com/fajarardiyanto/chat-application/internal/mapper"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ChatHandler struct {
	conversationRepository repository.ConversationRepository
	messageRepository      repository.MessageRepository
}

func NewChatHandler(
	conversationRepository repository.ConversationRepository,
	messageRepository repository.MessageRepository,
) *ChatHandler {
	return &ChatHandler{
		conversationRepository: conversationRepository,
		messageRepository:      messageRepository,
	}
}

func (s *ChatHandler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	conversationId := utils.QueryParam(r, "conversationId")

	if !commons.IsAllowedToChat(r) {
		config.GetLogger().Error(exception.NotAllowedToSetPassword)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowedToSetPassword)
		return
	}

	req := &request.MessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err = s.conversationRepository.FindByConversationId(conversationId); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, exception.ConversationNotFound)
		return
	}

	data := model.Message{
		Content:        req.Content,
		MessageType:    constant.MessageType[req.ContentType],
		SenderId:       token.UserId,
		CreatedAt:      time.Now(),
		ConversationId: conversationId,
		Uuid:           uuid.NewString(),
	}

	if len(req.Documents) > 0 {
		data.DocumentAttached = true
	}

	if err = s.messageRepository.CreateMessage(data); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OnMsg(config.GetConfig().Message, data)

	model.MessageSuccess(w, http.StatusOK, mapper.MessageMapper(data))
}
