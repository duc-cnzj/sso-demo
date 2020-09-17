package models

import (
	"time"
)

type Role struct {
	ID uint `gorm:"primary_key" json:"id"`

	Text string `gorm:"type:varchar(100);not null;" json:"text"`
	Name string `gorm:"type:varchar(50);unique_index;not null;" json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Permissions []Permission `gorm:"many2many:role_permission;" json:"permissions"`
	Users       []User       `gorm:"many2many:user_role;" json:"users"`
}
