package response

import "time"

type MessageResponse struct {
	ConversationId string `json:"conversation_id"`
	Content        string `json:"content"`
}

type MessageWsResponse struct {
	Event        string                      `json:"event"`
	Conversation ConversationMessageResponse `json:"conversation"`
	Data         InfoMessageResponse         `json:"data"`
}

type AttachmentMessageResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type UserMessageResponse struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ConversationMessageResponse struct {
	Agent          UserMessageResponse `json:"agent,omitempty"`
	Contact        UserMessageResponse `json:"contact,omitempty"`
	ConversationId string              `json:"conversation_id"`
}

type InfoMessageResponse struct {
	MessageId  string                      `json:"message_id"`
	Content    string                      `json:"content"`
	Attachment []AttachmentMessageResponse `json:"attachment,omitempty"`
	Timestamp  time.Time                   `json:"timestamp"`
	SenderType string                      `json:"sender_type"`
}
