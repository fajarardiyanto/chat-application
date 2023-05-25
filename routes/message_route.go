package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/controller"
	"github.com/fajarardiyanto/chat-application/internal/middleware"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func ChatRoute(r *mux.Router) {
	conversationService := services.NewConversationService()
	messageService := services.NewMessageService()

	chatHandler := controller.NewChatHandler(conversationService, messageService)

	secure := r.PathPrefix("/v1").Subrouter()
	secure.Use(middleware.AuthMiddleware)

	// /account/:accountId
	secure.HandleFunc("/conversation/{conversationId}/message", chatHandler.SendMessageHandler).Methods(http.MethodPost)
}
