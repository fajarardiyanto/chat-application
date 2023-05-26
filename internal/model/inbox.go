package model

type Inbox struct {
	Id                   int32  `gorm:"column:id"`
	CreatedAt            string `gorm:"column:created_at"`
	UpdatedAt            string `gorm:"column:updated_at"`
	ChannelId            string `gorm:"column:channel_id"`
	ChannelType          string `gorm:"column:channel_type"`
	EnableAutoAssignment bool   `gorm:"column:enable_auto_assignment"`
	GreetingEnabled      bool   `gorm:"column:greeting_enabled"`
	GreetingMessage      string `gorm:"column:greeting_message"`
	Name                 string `gorm:"column:name"`
	OooMessage           string `gorm:"column:ooo_message"`
	Timezone             string `gorm:"column:timezone"`
	Uuid                 string `gorm:"column:uuid"`
	WorkingHoursEnabled  bool   `gorm:"column:working_hours_enabled"`
	AccountId            string `gorm:"column:account_id"`
	Target               int32  `gorm:"column:target"`
}

func (*Inbox) TableName() string {
	return "inboxes"
}
