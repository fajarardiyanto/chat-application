package listener

import (
	"fmt"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/gorilla/websocket"
	"log"
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
	contactRepository      repository.ContactRepository
}

func NewEventListener(
	agentProfileRepository repository.AgentProfileRepository,
	contactRepository repository.ContactRepository,
) *EventListener {
	return &EventListener{
		agentProfileRepository: agentProfileRepository,
		contactRepository:      contactRepository,
	}
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

					var senderType string
					agentResponse := new(response.UserMessageResponse)
					contactResponse := new(response.UserMessageResponse)

					if msg.SenderType == 1 {
						agent, err := s.agentProfileRepository.FindAgentProfileById(msg.SenderId)
						if err != nil {
							config.GetLogger().Error(err.Error())
						}

						agentResponse = &response.UserMessageResponse{
							Name: fmt.Sprintf("%s %s", agent.FirstName, agent.LastName),
							Id:   msg.SenderId,
						}

						senderType = "AGENT"
					}

					if msg.SenderType == 0 {
						contact, err := s.contactRepository.FindById(msg.SenderId)
						if err != nil {
							config.GetLogger().Error(err.Error())
						}

						contactResponse = &response.UserMessageResponse{
							Name: contact.Name,
							Id:   msg.SenderId,
						}

						senderType = "CONTACT"
					}

					for client := range AgentClients {
						message := response.MessageWsResponse{
							Event: "MESSAGE_CREATED",
							Conversation: response.ConversationMessageResponse{
								ConversationId: msg.ConversationId,
							},
							Data: response.InfoMessageResponse{
								MessageId:  msg.Uuid,
								Content:    msg.Content,
								Timestamp:  msg.CreatedAt,
								SenderType: senderType,
							},
						}

						if agentResponse != nil {
							message.Conversation.Agent = *agentResponse
						}

						if contactResponse != nil {
							message.Conversation.Contact = *contactResponse
						}

						log.Println(client.ConversationId, "== ", msg.ConversationId)

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
								ConversationId: msg.ConversationId,
							},
							Data: response.InfoMessageResponse{
								MessageId:  msg.Uuid,
								Content:    msg.Content,
								Timestamp:  msg.CreatedAt,
								SenderType: senderType,
							},
						}

						if agentResponse != nil {
							message.Conversation.Agent = *agentResponse
						}

						if contactResponse != nil {
							message.Conversation.Contact = *contactResponse
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
