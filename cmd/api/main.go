package api

import (
	"fmt"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/handler"
	"github.com/fajarardiyanto/chat-application/internal/middleware"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"net/http"
)

var CmdAPI = &cobra.Command{
	Use:   "api",
	Short: "Start api server",
	RunE:  Api,
}

func Api(cmd *cobra.Command, args []string) error {
	config.Database(config.GetConfig().Database.SQL)

	userSvc := services.NewUserService()
	chatSvc := services.NewChatService()
	wsHandler := services.NewWSHandler()

	userHandler := handler.NewUserHandler(userSvc)
	chatHandler := handler.NewChatHandler(chatSvc)

	r := mux.NewRouter()
	r.Use(middleware.SetMiddlewareJSON)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		utils.MessageSuccess(w, http.StatusOK, "OK")
	}).Methods(http.MethodGet)

	r.HandleFunc("/register", userHandler.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/login", userHandler.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/contact-list", userHandler.ContactListHandler).Methods(http.MethodGet)
	r.HandleFunc("/update/status", userHandler.UpdateStatusHandler).Methods(http.MethodPost)

	r.HandleFunc("/create-chat", chatHandler.CreateMessageHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat-history", chatHandler.ChatHistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/files", chatHandler.SaveFileChat).Methods(http.MethodPost)
	r.HandleFunc("/static", chatHandler.StaticFile).Methods(http.MethodGet)

	go wsHandler.BroadcastWebSocket()
	r.HandleFunc("/ws", wsHandler.ServeWs)

	port := fmt.Sprintf(":%s", config.GetConfig().Port)
	config.GetLogger().Success("http server is starting on %s", port)

	hand := cors.Default().Handler(r)
	if err := http.ListenAndServe(port, hand); err != nil {
		config.GetLogger().Error(err)
		return err
	}

	return nil
}
