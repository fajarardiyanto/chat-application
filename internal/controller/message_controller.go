package controller

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/common/constant"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/common/mapper"
	"github.com/fajarardiyanto/chat-application/internal/common/validation"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/request"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type MessageController struct {
	conversationRepository repository.ConversationRepository
	messageRepository      repository.MessageRepository
	accountRepository      repository.AccountRepository
}

func NewChatHandler(
	conversationRepository repository.ConversationRepository,
	messageRepository repository.MessageRepository,
	accountRepository repository.AccountRepository,
) *MessageController {
	return &MessageController{
		conversationRepository: conversationRepository,
		messageRepository:      messageRepository,
		accountRepository:      accountRepository,
	}
}

func (s *MessageController) AgentSendMessageHandler(w http.ResponseWriter, r *http.Request) {
	conversationId := utils.QueryParam(r, "conversationId")
	accountId := utils.QueryParam(r, "accountId")

	req := &request.MessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		config.GetLogger().Error(err)
		model.MessageError(w, http.StatusBadRequest, exception.ErrorDecodeRequest)
		return
	}

	if _, err := s.accountRepository.FindAccountByAccountId(accountId); err != nil {
		config.GetLogger().Error(exception.AccountNotFound)
		model.MessageError(w, http.StatusNotFound, exception.AccountNotFound)
		return
	}

	if !validation.IsAllowedToChat(r) {
		config.GetLogger().Error(exception.NotAllowedToSetPassword)
		model.MessageError(w, http.StatusUnauthorized, exception.NotAllowedToSetPassword)
		return
	}

	token, err := auth.ExtractTokenAgent(r)
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
		MessageType:    constant.MessageType(req.ContentType),
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
