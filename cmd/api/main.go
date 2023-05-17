package api

import (
	"fmt"
	"github.com/fajarardiyanto/chat-application/routes"
	"net/http"
	"os"
	"os/signal"

	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/middleware"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
)

var CmdAPI = &cobra.Command{
	Use:   "api",
	Short: "Start api server",
	RunE:  Api,
}

func Api(cmd *cobra.Command, args []string) error {
	config.Database(config.GetConfig().Database.SQL)

	wsHandler := services.NewWSHandler()

	r := mux.NewRouter()
	r.Use(middleware.SetMiddlewareJSON)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		model.MessageSuccess(w, http.StatusOK, "OK")
	}).Methods(http.MethodGet)

	routes.AgentRoute(r)

	secure := r.PathPrefix("/auth").Subrouter()
	secure.Use(middleware.AuthMiddleware)

	go wsHandler.BroadcastWebSocket()
	secure.HandleFunc("/ws", wsHandler.ServeWs)

	go func() {
		port := fmt.Sprintf(":%s", config.GetConfig().Port)
		config.GetLogger().Success("http server is starting on %s", port)

		hand := cors.Default().Handler(r)
		if err := http.ListenAndServe(port, hand); err != nil {
			config.GetLogger().Error(err).Quit()
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	config.GetLogger().Debug("Got Signal: %v", sig)

	os.Exit(1)

	return nil
}
