package handler

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"net/http"
)

type ChatHandler struct {
	repo repo.ChatRepository
}

func NewChatHandler(repo repo.ChatRepository) *ChatHandler {
	return &ChatHandler{repo: repo}
}

func (s *ChatHandler) CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	u := &model.MessageRequest{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		model.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	req := model.Chat{
		Msg:  u.Message,
		To:   u.To,
		From: token.ID,
	}

	chat, err := s.repo.CreateChat(req)
	if err != nil {
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OnMsg(config.GetConfig().Message, chat)

	model.MessageSuccess(w, http.StatusOK, chat)
}

func (s *ChatHandler) ChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	to := utils.QueryString(r, "to")

	token, err := auth.ExtractTokenID(r)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	chats, err := s.repo.GetChat(token.ID, to)
	if err != nil {
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	model.MessageSuccess(w, http.StatusOK, chats)
}

func (s *ChatHandler) SaveFileChat(w http.ResponseWriter, r *http.Request) {
	file, err := s.repo.SaveFileChat(r)
	if err != nil {
		model.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	model.MessageSuccess(w, http.StatusOK, file)
}

func (s *ChatHandler) StaticFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")

	file := r.URL.Query().Get("file")

	url := "external/files/" + file
	http.ServeFile(w, r, url)
}
