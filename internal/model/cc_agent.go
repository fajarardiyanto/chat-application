package model

import "time"

type CCAgent struct {
	AccountId   string    `gorm:"column:account_id"`
	AgentId     string    `gorm:"column:agent_id"`
	AgentStatus int32     `gorm:"column:agent_status"`
	Available   bool      `gorm:"column:available"`
	LastSeenAt  time.Time `gorm:"column:last_seen_at"`
	PubSubToken string    `gorm:"column:pubsub_token"`
	Audit
}

func (*CCAgent) TableName() string {
	return "cc_agent"
}
