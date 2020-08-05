package models

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"sso/config/env"
	"sso/utils/helper"
	"time"
)

type User struct {
	ID          uint           `gorm:"primary_key" json:"id"`
	UserName    string         `gorm:"type:varchar(255);" json:"user_name"`
	Email       string         `gorm:"type:varchar(100);unique_index" json:"email"`
	AccessToken sql.NullString `gorm:"type:varchar(255);index" json:"-"`
	Password    string         `gorm:"type:varchar(255);" json:"-"`
	LastLoginAt time.Time      `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `sql:"index" json:"deleted_at"`
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

func (User) FindById(id uint, env *env.Env) *User {
	user := &User{}

	err := env.GetDB().Where("id = ?", id).First(user)
	if err.Error != nil {
		log.Println("findById", err)
		return nil
	}

	return user
}

func (User) TableName() string {
	return "users"
}

func (user *User) GetJwtID() string {
	return user.Email
}

func (user *User) GenerateAccessToken(env *env.Env) string {
	var try int
	for {
		if try > 10 {
			panic("error GenerateAccessToken try > 10")
		}

		str := helper.RandomString(64)
		AccessToken := sql.NullString{
			String: str,
			Valid:  true,
		}

		exists := env.GetDB().Table(User{}.TableName()).Where("access_token = ?", str).Find(&User{})
		if exists.Error != nil && errors.Is(gorm.ErrRecordNotFound, exists.Error){
			user.AccessToken = AccessToken
			env.GetDB().Model(user).Update("access_token", AccessToken)
			return str
		}

		try++
	}
}
