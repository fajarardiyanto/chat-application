package ws

import (
	"encoding/json"
	"fmt"
	"github.com/fajarardiyanto/chat-application/internal/common/exception"
	"github.com/fajarardiyanto/chat-application/internal/common/validation"
	"github.com/fajarardiyanto/chat-application/internal/repository"
	"github.com/fajarardiyanto/chat-application/internal/services/ws/listener"
	"net/http"
	"sync"
	"time"

	"github.com/fajarardiyanto/chat-application/config"
	"github.com/fajarardiyanto/chat-application/internal/model"
	"github.com/fajarardiyanto/chat-application/pkg/auth"
)

type wsHandler struct {
	sync.Mutex
	ccAgentRepository      repository.CCAgentRepository
	conversationRepository repository.ConversationRepository
	agentProfileRepository repository.AgentProfileRepository
}

func NewWSHandler(
	ccAgentRepository repository.CCAgentRepository,
	conversationRepository repository.ConversationRepository,
	agentProfileRepository repository.AgentProfileRepository,
) *wsHandler {
	return &wsHandler{
		ccAgentRepository:      ccAgentRepository,
		conversationRepository: conversationRepository,
		agentProfileRepository: agentProfileRepository,
	}
}

func (s *wsHandler) ServeWsAgent(w http.ResponseWriter, r *http.Request) {
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		config.GetLogger().Error(err)
		return
	}

	client := &listener.AgentClient{Conn: ws}

	s.Lock()
	listener.AgentClients[client] = true
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
	delete(listener.AgentClients, client)
}

func (s *wsHandler) Receiver(client *listener.AgentClient) {
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

func (s *wsHandler) OnClose(client *listener.AgentClient) {
	if err := client.Conn.Close(); err != nil {
		return
	}
}

func (s *wsHandler) WriteMessage(client *listener.AgentClient, msg interface{}) {
	if err := client.Conn.WriteJSON(msg); err != nil {
		return
	}
}

func (s *wsHandler) Ping(client *listener.AgentClient) {
	defer func() {
		time.Sleep(20 * time.Second)
		go s.Ping(client)
	}()

	s.WriteMessage(client, map[string]interface{}{
		"event": "PING",
	})
}
