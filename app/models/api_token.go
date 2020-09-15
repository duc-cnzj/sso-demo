package models

import (
	"time"
)

type ApiToken struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	ApiToken  string     `gorm:"type:varchar(255);index" json:"api_token"`
	LastUseAt *time.Time `json:"last_use_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`

	UserID uint `json:"user_id"`
	User   User `json:"user"`
}
