package model

import "time"

type Message struct {
	Id               int32     `gorm:"column:id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	Content          string    `gorm:"column:content"`
	IsPrivate        bool      `gorm:"column:is_private"`
	MessageType      int32     `gorm:"column:message_type"`
	SenderId         string    `gorm:"column:sender_id"`
	SenderType       int       `gorm:"column:sender_type"`
	Uuid             string    `gorm:"column:uuid"`
	ConversationId   string    `gorm:"column:conversation_id"`
	Deleted          bool      `gorm:"column:deleted"`
	DocumentAttached bool      `gorm:"column:document_attached"`
}

func (*Message) TableName() string {
	return "messages"
}
