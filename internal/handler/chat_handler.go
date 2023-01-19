package handler

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repo"
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
		utils.MessageError(w, http.StatusBadRequest, "error decoding request object")
		return
	}

	req := model.Chat{
		Msg: u.Message,
		To:  u.To,
	}

	chat, err := s.repo.CreateChat(req)
	if err != nil {
		utils.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OnMsg(config.GetConfig().Message, chat)

	utils.MessageSuccess(w, http.StatusOK, chat)
}

func (s *ChatHandler) ChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	from := utils.QueryString(r, "from")
	to := utils.QueryString(r, "to")

	chats, err := s.repo.GetChat(from, to)
	if err != nil {
		utils.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.MessageSuccess(w, http.StatusOK, chats)
}

func (s *ChatHandler) SaveFileChat(w http.ResponseWriter, r *http.Request) {
	file, err := s.repo.SaveFileChat(r)
	if err != nil {
		utils.MessageError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.MessageSuccess(w, http.StatusOK, file)
}

func (s *ChatHandler) StaticFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")

	file := r.URL.Query().Get("file")

	url := "external/files/" + file
	http.ServeFile(w, r, url)
}
