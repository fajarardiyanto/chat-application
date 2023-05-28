package ws

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"log"
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
	ConversationId string
	InboxId        string
}

type webWidgetWsHandler struct {
	sync.Mutex
	inboxRepository        repository.InboxRepository
	contactInboxRepository repository.ContactInboxRepository
	conversationRepository repository.ConversationRepository
	contactRepository      repository.ContactRepository
}

func NewWebWidgetWSHandler(
	inboxRepository repository.InboxRepository,
	contactInboxRepository repository.ContactInboxRepository,
	conversationRepository repository.ConversationRepository,
	contactRepository repository.ContactRepository,
) *webWidgetWsHandler {
	return &webWidgetWsHandler{
		inboxRepository:        inboxRepository,
		contactInboxRepository: contactInboxRepository,
		conversationRepository: conversationRepository,
		contactRepository:      contactRepository,
	}
}

func (s *webWidgetWsHandler) ServeWsAgent(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	client := &webWidgetClient{Conn: ws}

	websiteToken := utils.QueryString(r, "websiteToken")
	if websiteToken == "" {
		config.GetLogger().Error(exception.WebsiteTokenMissing)
		s.WriteMessage(client, map[string]interface{}{
			"message": exception.WebsiteTokenMissing,
		})
		s.OnClose(client)
		return
	}

	s.Lock()
	webWidgetClients[client] = true
	s.Unlock()

	token, err := auth.ExtractTokenContact(r)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	contactInbox, err := s.contactInboxRepository.FindBySourceId(token.SourceId)
	if err != nil {
		config.GetLogger().Error(err.Error())
		s.WriteMessage(client, map[string]interface{}{
			"message": err.Error(),
		})
		s.OnClose(client)
		return
	}

	contact, err := s.contactRepository.FindById(contactInbox.ContactId)
	if err != nil {
		config.GetLogger().Error(err.Error())
		s.WriteMessage(client, map[string]interface{}{
			"message": err.Error(),
		})
		s.OnClose(client)
		return
	}

	conversationId := uuid.NewString()
	conversation, err := s.conversationRepository.FindByContactInboxId(contactInbox.Uuid)
	if err != nil {
		conversationData := model.Conversation{
			Uuid:           conversationId,
			ContactInboxId: contactInbox.Uuid,
			ContactId:      contactInbox.ContactId,
			InboxId:        contactInbox.InboxId,
			Status:         "OPEN",
			AccountId:      contact.AccountId,
		}

		if err = s.conversationRepository.CreateConversation(conversationData); err != nil {
			config.GetLogger().Error(err.Error())
			s.WriteMessage(client, map[string]interface{}{
				"message": err.Error(),
			})
			s.OnClose(client)
			return
		}

		s.WriteMessage(client, map[string]interface{}{
			"message": conversationData,
		})
	} else {
		if conversation.Status != "OPEN" {
			config.GetLogger().Error(exception.ConversationClosed)
			s.WriteMessage(client, map[string]interface{}{
				"message": exception.ConversationClosed,
			})
			s.OnClose(client)
			return
		}
		conversationId = conversation.Uuid
	}

	s.Lock()
	client.ConversationId = conversationId
	s.Unlock()
	config.GetLogger().Info("%s is connected", conversationId)

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

					//contact, err := s.contactRepository.FindById(msg.SenderId)
					//if err != nil {
					//	config.GetLogger().Error(err.Error())
					//}

					for client := range webWidgetClients {
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

						log.Println("HEREEEE CONTACT", client.ConversationId, msg.ConversationId)
						if client.ConversationId == msg.ConversationId {
							if err = client.Conn.WriteJSON(message); err != nil {
								config.GetLogger().Error("%s is offline", client.InboxId)
								if err = client.Conn.Close(); err != nil {
									config.GetLogger().Error(err.Error())
									return
								}

								s.Lock()
								delete(webWidgetClients, client)
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

func (s *webWidgetWsHandler) CheckIfOutOfWorkingHours(inboxId string) *response.OperationalHoursResponse {
	return nil
}
