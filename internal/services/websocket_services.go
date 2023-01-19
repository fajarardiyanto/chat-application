package services

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var clients = make(map[*Client]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHandler struct {
	sync.Mutex
}

type Client struct {
	Conn     *websocket.Conn
	Username string
}

func NewWSHandler() *wsHandler {
	return &wsHandler{}
}

func (s *wsHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.GetLogger().Error(err)
	}

	client := &Client{Conn: ws}

	s.Lock()
	clients[client] = true
	s.Unlock()

	s.Lock()
	client.Username = utils.QueryString(r, "user")
	s.Unlock()

	s.Receiver(client)

	config.GetLogger().Error("exiting %s", ws.RemoteAddr().String())
	delete(clients, client)
}

func (s *wsHandler) Receiver(client *Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			config.GetLogger().Error(err)
			return
		}

		m := &model.Chat{}
		if err = json.Unmarshal(p, m); err != nil {
			config.GetLogger().Error("error while unmarshalling chat %v", err)
			continue
		}
	}
}

func (s *wsHandler) BroadcastWebSocket() {
	config.GetLogger().Info("Broadcaster started")

	s.OnMsg()
}

func (s *wsHandler) OnMsg() {
	go func() {
		config.GetRabbitMQ().Consumer(interfaces.RabbitMQOptions{
			Exchange: config.GetConfig().Message,
			NoWait:   true},
			func(m interfaces.Messages, cid interfaces.ConsumerCallbackIsDone) {
				var msg model.Chat
				if err := m.Decode(&msg); err == nil {
					config.GetLogger().Info("Message Receive %v", msg)
					for client := range clients {
						message := model.Message{
							Type: "MESSAGE_TEXT",
							Chat: msg,
						}
						if client.Username == msg.From || client.Username == msg.To {
							if err = client.Conn.WriteJSON(message); err != nil {
								config.GetLogger().Error("%s is offline", client.Username)
								if err = client.Conn.Close(); err != nil {
									config.GetLogger().Error(err.Error())
									return
								}
								s.Lock()
								delete(clients, client)
								s.Unlock()
							}
						}
					}
				}
			})
	}()
}
