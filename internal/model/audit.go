package model

import "time"

type Audit struct {
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedBy string    `json:"created_by" gorm:"column:created_by"`
	UpdatedBy string    `json:"updated_by" gorm:"column:updated_by"`
}
