package routes

import (
	"github.com/fajarardiyanto/chat-application/internal/controller"
	"github.com/fajarardiyanto/chat-application/internal/middleware"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func AgentRoute(r *mux.Router) {
	agentProfileService := services.NewAgentProfileService()
	accountService := services.NewAccountService()
	agentCredential := services.NewAgentCredentialService()

	agentProfileController := controller.NewAgentProfileController(agentProfileService, accountService, agentCredential)

	agentProfile := r.PathPrefix("/v1").Subrouter()
	agentProfile.HandleFunc("/agent", agentProfileController.RegisterHandler).Methods(http.MethodPost)
	agentProfile.HandleFunc("/agent/login", agentProfileController.Login).Methods(http.MethodPost)

	agentProfileSecure := r.PathPrefix("/v1").Subrouter()
	agentProfileSecure.Use(middleware.AuthMiddleware)
	agentProfileSecure.HandleFunc("/agent/password", agentProfileController.SetPassword).Methods(http.MethodPost)
	agentProfileSecure.HandleFunc("/agent", agentProfileController.GetAgent).Methods(http.MethodGet)
	agentProfileSecure.HandleFunc("/agent/account/{accountId}/agents", agentProfileController.GetAllAgentByAccountId).Methods(http.MethodGet)
	agentProfileSecure.HandleFunc("/agent/{agentId}", agentProfileController.UpdateAgentProfileById).Methods(http.MethodPut)
	agentProfileSecure.HandleFunc("/agent/{agentId}", agentProfileController.DeleteAgentProfileById).Methods(http.MethodDelete)
}
