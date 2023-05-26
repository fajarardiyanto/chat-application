package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/controller"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func ContactRoute(r *mux.Router) {
	contactService := services.NewContactService()
	contactInboxService := services.NewContactInboxService()
	channelWebWidgetService := services.NewChannelWebWidgetService()
	inboxService := services.NewInboxService()

	contactController := controller.NewContactController(contactService, contactInboxService, channelWebWidgetService, inboxService)

	secure := r.PathPrefix("/v1/widget").Subrouter()
	//secure.Use(middleware.AuthMiddleware)

	secure.HandleFunc("/website-token/{websiteToken}/auth", contactController.RegisterContact).Methods(http.MethodPost)
}
