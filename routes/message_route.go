package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/controller"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/fajarardiyanto/chat-application/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func MessageRoute(r *mux.Router) {
	conversationService := services.NewConversationService()
	messageService := services.NewMessageService()
	accountService := services.NewAccountService()

	chatHandler := controller.NewChatHandler(conversationService, messageService, accountService)

	secure := r.PathPrefix("/v1").Subrouter()
	secure.Use(middleware.AuthMiddleware)

	secure.HandleFunc("/account/{accountId}/conversation/{conversationId}/message", chatHandler.AgentSendMessageHandler).Methods(http.MethodPost)
}
