package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/fajarardiyanto/chat-application/internal/services/ws"
	"github.com/fajarardiyanto/chat-application/middleware"
	"github.com/gorilla/mux"
)

func WsRoute(r *mux.Router) {
	ccAgentService := services.NewCCAgentService()
	conversationService := services.NewConversationService()
	agentProfileService := services.NewAgentProfileService()

	wsHandler := ws.NewWSHandler(ccAgentService, conversationService, agentProfileService)

	secure := r.PathPrefix("/ws").Subrouter()
	secure.Use(middleware.AuthMiddleware)

	go wsHandler.BroadcastWebSocket()
	secure.HandleFunc("/cable", wsHandler.ServeWsAgent)
}
