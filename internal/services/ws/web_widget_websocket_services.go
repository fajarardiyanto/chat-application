package ws

import (
	"encoding/json"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/controller/dto/response"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/internal/services/ws/listener"
	"github.com/fajarardiyanto/chat-application/pkg/utils"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"

	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
)

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

func (s *webWidgetWsHandler) ServeWsWebWidget(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	client := &listener.WebWidgetClient{Conn: ws}

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
	listener.WebWidgetClients[client] = true
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
	delete(listener.WebWidgetClients, client)
}

func (s *webWidgetWsHandler) Receiver(client *listener.WebWidgetClient) {
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

func (s *webWidgetWsHandler) OnClose(client *listener.WebWidgetClient) {
	if err := client.Conn.Close(); err != nil {
		return
	}
}

func (s *webWidgetWsHandler) WriteMessage(client *listener.WebWidgetClient, msg interface{}) {
	if err := client.Conn.WriteJSON(msg); err != nil {
		return
	}
}

func (s *webWidgetWsHandler) Ping(client *listener.WebWidgetClient) {
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
