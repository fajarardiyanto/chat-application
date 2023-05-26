package ws

import (
	"encoding/json"
	"fmt"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/common/validation"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"net/http"
	"sync"
	"time"

	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
	"github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/gorilla/websocket"
)

var webWidgetClients = make(map[*webWidgetClient]bool)

type webWidgetClient struct {
	Conn           *websocket.Conn
	Username       string
	ConversationId string
}

type webWidgetWsHandler struct {
	sync.Mutex
	ccAgentRepository      repository.CCAgentRepository
	conversationRepository repository.ConversationRepository
	agentProfileRepository repository.AgentProfileRepository
}

func NewWebWidgetWSHandler(
	ccAgentRepository repository.CCAgentRepository,
	conversationRepository repository.ConversationRepository,
	agentProfileRepository repository.AgentProfileRepository,
) *webWidgetWsHandler {
	return &webWidgetWsHandler{
		ccAgentRepository:      ccAgentRepository,
		conversationRepository: conversationRepository,
		agentProfileRepository: agentProfileRepository,
	}
}

func (s *webWidgetWsHandler) ServeWsAgent(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	client := &webWidgetClient{Conn: ws}

	s.Lock()
	webWidgetClients[client] = true
	s.Unlock()

	token, err := auth.ExtractTokenAgent(r)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	s.Lock()
	client.Username = token.UserId
	s.Unlock()
	config.GetLogger().Info("%s is connected", token.UserId)

	if !validation.IsAllowedToChat(r) {
		config.GetLogger().Error(exception.NotAllowedToChat)
		s.WriteMessage(client, map[string]interface{}{
			"message": "This User with role " + token.Role + " not allowed to chat with customer",
		})
		s.OnClose(client)
		return
	}

	ccAgent, err := s.ccAgentRepository.FindCCAgentByAgentId(token.UserId)
	if err != nil {
		errMsg := fmt.Sprintf("%s agent not present in chat platform", token.UserId)
		config.GetLogger().Error(errMsg)
		s.WriteMessage(client, map[string]interface{}{
			"message": errMsg,
		})
		s.OnClose(client)
		return
	}

	if token.AccountId != ccAgent.AccountId {
		errMsg := "Looks like agent is not in the same account. Check with your manager"
		config.GetLogger().Error(errMsg)
		s.WriteMessage(client, map[string]interface{}{
			"message": errMsg,
		})
		s.OnClose(client)
		return
	}

	if conversation, err := s.conversationRepository.FindByAgentId(ccAgent.AgentId); err == nil {
		s.Lock()
		client.ConversationId = conversation.Uuid
		s.Unlock()
	}

	s.Ping(client)
	s.Receiver(client)

	config.GetLogger().Error("exiting %s", ws.RemoteAddr().String())
	delete(webWidgetClients, client)
}

func (s *webWidgetWsHandler) Receiver(client *webWidgetClient) {
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

func (s *webWidgetWsHandler) BroadcastWebSocket() {
	config.GetLogger().Info("Broadcaster started")

	s.OnMsg()
}

func (s *webWidgetWsHandler) OnMsg() {
	go func() {

		config.GetRabbitMQ().Consumer(interfaces.RabbitMQOptions{
			Exchange: config.GetConfig().Message,
			NoWait:   true},
			func(m interfaces.Messages, cid interfaces.ConsumerCallbackIsDone) {
				var msg model.Message
				if err := m.Decode(&msg); err == nil {
					config.GetLogger().Info("Message Receive %v", msg)

					agent, err := s.agentProfileRepository.FindAgentProfileById(msg.SenderId)
					if err != nil {
						config.GetLogger().Error(err.Error())
					}

					for client := range agentClients {
						message := response.MessageWsResponse{
							Event: "MESSAGE_CREATED",
							Conversation: response.ConversationMessageResponse{
								Agent: response.UserMessageResponse{
									Name: msg.SenderId,
									Id:   fmt.Sprintf("%s %s", agent.FirstName, agent.LastName),
								},
								ConversationId: msg.ConversationId,
							},
							Data: response.InfoMessageResponse{
								MessageId:  msg.Uuid,
								Content:    msg.Content,
								Timestamp:  msg.CreatedAt,
								SenderType: "AGENT",
							},
						}

						if client.ConversationId == msg.ConversationId && client.Username == msg.SenderId {
							if err = client.Conn.WriteJSON(message); err != nil {
								config.GetLogger().Error("%s is offline", client.Username)
								if err = client.Conn.Close(); err != nil {
									config.GetLogger().Error(err.Error())
									return
								}

								s.Lock()
								delete(agentClients, client)
								s.Unlock()
							}
						}
					}
				}
			})
	}()
}

func (s *webWidgetWsHandler) OnClose(client *webWidgetClient) {
	if err := client.Conn.Close(); err != nil {
		return
	}
}

func (s *webWidgetWsHandler) WriteMessage(client *webWidgetClient, msg interface{}) {
	if err := client.Conn.WriteJSON(msg); err != nil {
		return
	}
}

func (s *webWidgetWsHandler) Ping(client *webWidgetClient) {
	defer func() {
		time.Sleep(20 * time.Second)
		go s.Ping(client)
	}()

	s.WriteMessage(client, map[string]interface{}{
		"event": "PING",
	})
}
