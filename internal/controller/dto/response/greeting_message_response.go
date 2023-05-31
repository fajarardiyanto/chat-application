package response

type GreetingMessageResponse struct {
	Conversation ConversationGreetingMessage `json:"conversation"`
}

type ConversationGreetingMessage struct {
	ConversationID string `json:"conversation_id"`
	ContactInboxId string `json:"contact_inbox_id"`
	ContactId      string `json:"contact_id"`
	InboxId        string `json:"inbox_id"`
}
