package models

import (
	"log"
	"sso/config/env"
	"time"
)

type User struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	UserName    string     `gorm:"type:varchar(255);" json:"user_name"`
	Email       string     `gorm:"type:varchar(100);unique_index" json:"email"`
	Password    string     `gorm:"type:varchar(255);" json:"password"`
	LastLoginAt time.Time  `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}

func (User) FindByEmail(email string, env *env.Env) *User {
	user := &User{}

	err := env.GetDB().Where("email = ?", email).First(user)
	if err.Error != nil {
		log.Println("findByEmail", err)
		return nil
	}

	return user
}

func (User) TableName() string {
	return "users"
}
