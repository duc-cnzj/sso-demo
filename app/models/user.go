package models

import (
	"time"
)

type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	UserName    string     `gorm:"type:varchar(255);not null" json:"user_name"`
	Email       string     `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	LogoutToken string     `gorm:"type:varchar(255);index;default:'';not null" json:"-"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`

	Permissions []Permission `gorm:"many2many:user_permission;" json:"permissions"`
	Roles       []Role       `gorm:"many2many:user_role;" json:"roles"`
	ApiTokens   []ApiToken   `json:"api_tokens"`

	CurrentRoles       []string `gorm:"-" json:"-"`
	CurrentPermissions []string `gorm:"-" json:"-"`
}

func (User) TableName() string {
	return "users"
}
