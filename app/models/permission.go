package models

import (
	"time"
)

type Permission struct {
	ID uint `gorm:"primary_key" json:"id"`

	Text    string `gorm:"type:varchar(100);not null;" json:"text"`
	Name    string `gorm:"type:varchar(50);not null;" json:"name"`
	Project string `json:"project"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Roles []Role `gorm:"many2many:role_permission;" json:"roles"`
	Users []User `gorm:"many2many:user_permission;" json:"users"`
}
