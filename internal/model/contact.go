package model

import "time"

type Contact struct {
	Id               int32     `gorm:"column:id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	Email            string    `gorm:"column:email"`
	Name             string    `gorm:"column:name"`
	Phone            string    `gorm:"column:phone"`
	Uuid             string    `gorm:"column:uuid"`
	AccountId        string    `gorm:"column:account_id"`
	DateOfBirth      string    `gorm:"column:date_of_birth"`
	Gender           string    `gorm:"column:gender"`
	MotherMaidenName string    `gorm:"column:mother_maiden_name"`
	Note             string    `gorm:"column:note"`
}

func (*Contact) TableName() string {
	return "contacts"
}
