package model

import "time"

type Chat struct {
	ID          string    `json:"id" gorm:"column:id"`
	From        string    `json:"from" gorm:"column:from_user"`
	To          string    `json:"to" gorm:"column:to_user"`
	Msg         string    `json:"message" gorm:"column:message"`
	MessageType int32     `json:"message_type" gorm:"column:message_type"`
	File        string    `json:"file,omitempty" gorm:"column:file"`
	FileSuffix  string    `json:"file_suffix,omitempty" gorm:"column:file_suffix"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (*Chat) TableName() string {
	return "chat"
}
