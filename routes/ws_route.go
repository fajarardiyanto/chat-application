package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/middleware"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/gorilla/mux"
)

func WsRoute(r *mux.Router) {
	ccAgentService := services.NewCCAgentService()
	conversationService := services.NewConversationService()
	agentProfileService := services.NewAgentProfileService()

	wsHandler := services.NewWSHandler(ccAgentService, conversationService, agentProfileService)

	secure := r.PathPrefix("/ws").Subrouter()
	secure.Use(middleware.AuthMiddleware)

	go wsHandler.BroadcastWebSocket()
	secure.HandleFunc("/cable", wsHandler.ServeWsAgent)
}
