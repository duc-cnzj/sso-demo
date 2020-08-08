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
	ID                uint           `gorm:"primary_key" json:"id"`
	UserName          string         `gorm:"type:varchar(255);not null" json:"user_name"`
	Email             string         `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	ApiToken          sql.NullString `gorm:"type:varchar(255);index" json:"-"`
	ApiTokenCreatedAt *time.Time     `json:"-"`
	LogoutToken       string         `gorm:"type:varchar(255);index;default:'';not null" json:"-"`
	Password          string         `gorm:"type:varchar(255);not null" json:"-"`
	LastLoginAt       *time.Time     `json:"last_login_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         *time.Time     `sql:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
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

func (User) FindByToken(token string, env *env.Env) *User {
	user := &User{}

	err := env.GetDB().Where("api_token = ?", token).First(user)
	if err.Error != nil {
		log.Println("FindByToken", err)
		return nil
	}

	return user
}

func (user *User) TokenExpired(env *env.Env) bool {
	seconds := time.Second * time.Duration(env.Config().SessionLifetime)
	if user.ApiToken.Valid &&
		user.ApiToken.String != "" &&
		time.Now().Before(user.ApiTokenCreatedAt.Add(seconds)) {
		return false
	}

	return true
}

func (user *User) GenerateAccessToken(env *env.Env) string {
	var (
		try   int
		str   string
		err   error
		reply interface{}
	)

	conn := env.RedisPool().Get()
	defer conn.Close()

	for {
		if try > 10 {
			panic("error GenerateAccessToken try > 10")
		}

		str = helper.RandomString(64)

		reply, err = conn.Do("GET", str)
		if err == nil && reply == nil {
			if env.Config().AccessTokenLifetime > 0 {
				reply, err = conn.Do("SETEX", str, env.Config().AccessTokenLifetime, user.ID)
			} else {
				reply, err = conn.Do("SET", str, user.ID)
			}
			log.Println(reply, err)
			if err == nil {
				return str
			}
		}

		try++
	}
}

func (user *User) GenerateApiToken(env *env.Env, forceFill bool) string {
	var try int

	log.Println(user)
	// 如果生成过token，并且没过期，则直接返回
	if !forceFill &&
		user.ApiToken.Valid &&
		user.ApiToken.String != "" &&
		!user.TokenExpired(env) {
		return user.ApiToken.String
	}

	for {
		if try > 10 {
			panic("error GenerateAccessToken try > 10")
		}

		str := helper.RandomString(64)
		AccessToken := sql.NullString{
			String: str,
			Valid:  true,
		}

		exists := env.GetDB().Table(User{}.TableName()).Where("api_token = ?", str).Find(&User{})
		if exists.Error != nil && errors.Is(gorm.ErrRecordNotFound, exists.Error) {
			user.ApiToken = AccessToken
			env.GetDB().Model(user).Update("api_token", AccessToken)
			env.GetDB().Model(user).Update(map[string]interface{}{"api_token": AccessToken, "api_token_created_at": time.Now()})
			return str
		}

		try++
	}
}

func (user *User) UpdateLastLoginAt(env *env.Env) {
	env.GetDB().Model(user).Update("last_login_at", time.Now())
}

func (user *User) GenerateLogoutToken(env *env.Env) {
	str := helper.RandomString(64)
	env.GetDB().Model(user).Update("logout_token", str)
}
