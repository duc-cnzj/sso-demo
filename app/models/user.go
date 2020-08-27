package models

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
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

	Permissions []Permission `gorm:"many2many:user_permission;" json:"permissions"`
	Roles       []Role       `gorm:"many2many:user_role;" json:"roles"`
}

func (User) TableName() string {
	return "users"
}

func (User) FindWithRoles(id int, env *env.Env) *User {
	user := &User{}

	if err := env.GetDB().Preload("Roles").Where("id = ?", id).First(user).Error; err != nil {
		log.Debug().Err(err).Msg("FindWithRoles")
	}

	return user
}
func (User) FindByEmail(email string, env *env.Env, wheres ...interface{}) *User {
	user := &User{}

	if err := env.GetDB().Where("email = ?", email).First(user, wheres...).Error; err != nil {
		log.Debug().Err(err).Msg("findByEmail")
		return nil
	}

	return user
}

func (User) FindById(id uint, env *env.Env) *User {
	user := &User{}

	if err := env.GetDB().Where("id = ?", id).First(user).Error; err != nil {
		log.Debug().Err(err).Msg("FindById")

		return nil
	}

	return user
}

func (User) FindByToken(token string, env *env.Env) *User {
	user := &User{}

	if err := env.GetDB().Where("api_token = ?", token).First(user).Error; err != nil {
		log.Debug().Err(err).Msg("FindByToken")
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
			log.Debug().Err(err).Interface("reply", reply).Msg("GenerateAccessToken")
			if err == nil {
				return str
			}
		}

		try++
	}
}

func (user *User) GenerateApiToken(env *env.Env, forceFill bool) string {
	var try int

	log.Debug().Interface("user", user).Msg("GenerateApiToken")
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

func (user *User) SyncRoles(roles []*Role, env *env.Env) error {
	return env.DBTransaction(func(tx *gorm.DB) error {
		if tx.Model(user).Association("Roles").Clear().Error != nil {
			return tx.Model(user).Association("Roles").Clear().Error
		}

		tx.Model(user).Association("Roles").Append(toRoleInterfaceSlice(roles)...)

		return nil
	})
}

func (user *User) SyncPermissions(permissions []interface{}, env *env.Env) error {
	return env.DBTransaction(func(tx *gorm.DB) error {
		if tx.Model(user).Association("Permissions").Clear().Error != nil {
			return tx.Model(user).Association("Permissions").Clear().Error
		}

		tx.Model(user).Association("Permissions").Append(permissions...)

		return nil
	})
}

// 登出用户，为了保证一处登出，处处登出，必须重置 api_token 和 logout_token
// api_token 重置之后会重定向到 sso login page, 但是 sso 依然是登陆状态
// 所以 sso 会重新生成 api_token，导致登出没有效果
// 因此引入 logout_token 如果 sso session 中的 logout_token 不一样，那么代表强制登出
// 以此来保证 用户登出时，api_token 失效，并且 sso 也登陆过期
func (user *User) ForceLogout(env *env.Env) {
	user.GenerateApiToken(env, true)
	user.GenerateLogoutToken(env)
}

func (user User) GeneratePwd(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

func toRoleInterfaceSlice(slice interface{}) []interface{} {
	roles := slice.([]*Role)
	newS := make([]interface{}, len(roles))
	for i, v := range roles {
		newS[i] = v
	}

	return newS
}
