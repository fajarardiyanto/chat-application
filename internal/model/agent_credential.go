package model

type AgentCredential struct {
	Audit
	Password string `gorm:"column:encrypted_password"`
	SaltKey  string `gorm:"column:salt_key"`
	UserName string `gorm:"column:username"`
}

func (*AgentCredential) TableName() string {
	return "agents_credentials"
}
