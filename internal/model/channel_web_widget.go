package model

type ChannelWebWidget struct {
	Id           int32  `gorm:"column:id"`
	CreatedAt    string `gorm:"column:created_at"`
	UpdatedAt    string `gorm:"column:updated_at"`
	Uuid         string `gorm:"column:uuid"`
	WebsiteToken string `gorm:"column:website_token"`
	WebsiteUrl   string `gorm:"column:website_url"`
	WelcomeTitle string `gorm:"column:welcome_title"`
	WidgetColor  string `gorm:"column:widget_color"`
	AccountId    string `gorm:"column:account_id"`
}

func (*ChannelWebWidget) TableName() string {
	return "channel_web_widget"
}
