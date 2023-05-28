package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/fajarardiyanto/chat-application/internal/services/ws"
	"github.com/fajarardiyanto/chat-application/middleware"
	"github.com/gorilla/mux"
)

func WsRoute(r *mux.Router) {
	//ccAgentService := services.NewCCAgentService()
	conversationService := services.NewConversationService()
	//agentProfileService := services.NewAgentProfileService()
	inboxService := services.NewInboxService()
	contactInboxService := services.NewContactInboxService()
	contactService := services.NewContactService()

	//agentWsHandler := ws.NewAgentWSHandler(ccAgentService, conversationService, agentProfileService)
	contactWsHandler := ws.NewWebWidgetWSHandler(inboxService, contactInboxService, conversationService, contactService)

	secure := r.PathPrefix("/ws").Subrouter()
	secure.Use(middleware.AuthMiddleware)

	//go agentWsHandler.BroadcastWebSocket()
	//secure.HandleFunc("/cable", agentWsHandler.ServeWsAgent)

	go contactWsHandler.BroadcastWebSocket()
	secure.HandleFunc("/web-widget/cable", contactWsHandler.ServeWsAgent)
}
