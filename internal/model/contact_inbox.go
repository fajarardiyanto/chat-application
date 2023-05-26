package model

import "time"

type ContactInbox struct {
	Id          int32     `gorm:"column:id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	PubSubToken string    `gorm:"column:pubsub_token"`
	SourceId    string    `gorm:"column:source_id"`
	Uuid        string    `gorm:"column:uuid"`
	ContactId   string    `gorm:"column:contact_id"`
	InboxId     string    `gorm:"column:inbox_id"`
}

func (*ContactInbox) TableName() string {
	return "contact_inbox"
}
