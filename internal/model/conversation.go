package model

import "time"

type Conversation struct {
	Id                  int32     `gorm:"column:id"`
	CreatedAt           time.Time `gorm:"column:created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at"`
	AgentId             string    `gorm:"column:agent_id"`
	State               int32     `gorm:"column:state"`
	Uuid                string    `gorm:"column:uuid"`
	AccountId           string    `gorm:"column:account_id"`
	ContactId           string    `gorm:"column:contact_id"`
	ContactInboxId      string    `gorm:"column:contact_inbox_id"`
	InboxId             string    `gorm:"column:inbox_id"`
	GreetingMessageSent bool      `gorm:"column:greeting_message_sent"`
	CommunicationType   string    `gorm:"column:communication_type"`
	LastSeenAgent       time.Time `gorm:"column:last_seen_agent"`
	EndTime             time.Time `gorm:"column:end_time"`
	StartTime           time.Time `gorm:"column:start_time"`
	FirstReplyAt        time.Time `gorm:"column:first_reply_at"`
	Status              string    `gorm:"column:status"`
	SummaryNoteNeeded   bool      `gorm:"column:summary_note_needed"`
}

func (*Conversation) TableName() string {
	return "conversation"
}
