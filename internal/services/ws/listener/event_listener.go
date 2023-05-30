package listener

import (
	"fmt"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/gorilla/websocket"
	"sync"
)

var AgentClients = make(map[*AgentClient]bool)
var WebWidgetClients = make(map[*WebWidgetClient]bool)

type AgentClient struct {
	Conn           *websocket.Conn
	Username       string
	ConversationId string
}

type WebWidgetClient struct {
	Conn           *websocket.Conn
	ConversationId string
	InboxId        string
}

type EventListener struct {
	sync.Mutex
	agentProfileRepository repository.AgentProfileRepository
}

func NewEventListener(agentProfileRepository repository.AgentProfileRepository) *EventListener {
	return &EventListener{agentProfileRepository: agentProfileRepository}
}

func (s *EventListener) OnMsg() {
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

					for client := range AgentClients {
						message := response.MessageWsResponse{
							Event: "MESSAGE_CREATED",
							Conversation: response.ConversationMessageResponse{
								Agent: response.UserMessageResponse{
									Name: fmt.Sprintf("%s %s", agent.FirstName, agent.LastName),
									Id:   msg.SenderId,
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
								delete(AgentClients, client)
								s.Unlock()
							}
						}
					}

					for client := range WebWidgetClients {
						message := response.MessageWsResponse{
							Event: "MESSAGE_CREATED",
							Conversation: response.ConversationMessageResponse{
								Agent: response.UserMessageResponse{
									Name: "test",
									Id:   msg.SenderId,
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

						if client.ConversationId == msg.ConversationId {
							if err = client.Conn.WriteJSON(message); err != nil {
								config.GetLogger().Error("%s is offline", client.InboxId)
								if err = client.Conn.Close(); err != nil {
									config.GetLogger().Error(err.Error())
									return
								}

								s.Lock()
								delete(WebWidgetClients, client)
								s.Unlock()
							}
						}
					}
				}
			})
	}()
}

func (s *EventListener) BroadcastWebSocket() {
	config.GetLogger().Info("Broadcaster started")

	s.OnMsg()
}
