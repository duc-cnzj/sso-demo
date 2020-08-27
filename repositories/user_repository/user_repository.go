package user_repository

import (
	"database/sql"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"sso/app/models"
	"sso/config/env"
	"sso/utils/helper"
	"time"
)

type UserRepository struct {
	env *env.Env
}

func NewUserRepository(env *env.Env) *UserRepository {
	return &UserRepository{
		env: env,
	}
}

func (repo *UserRepository) FindByEmail(email string, wheres ...interface{}) *models.User {
	user := &models.User{}

	if err := repo.env.GetDB().Where("email = ?", email).First(user, wheres...).Error; err != nil {
		log.Debug().Err(err).Msg("findByEmail")
		return nil
	}

	return user
}

func (repo *UserRepository) GeneratePwd(pwd string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func (repo *UserRepository) Create(user *models.User) error {
	if err := repo.env.GetDB().Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) FindById(id uint) (*models.User, error) {
	user := &models.User{}

	if err := repo.env.GetDB().Where("id = ?", id).First(user).Error; err != nil {
		log.Debug().Err(err).Msg("FindById")

		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) SyncRoles(user *models.User, roles []*models.Role) error {
	return repo.env.DBTransaction(func(tx *gorm.DB) error {
		if tx.Model(user).Association("Roles").Clear().Error != nil {
			return tx.Model(user).Association("Roles").Clear().Error
		}

		tx.Model(user).Association("Roles").Append(toRoleInterfaceSlice(roles)...)

		return nil
	})
}

// 登出用户，为了保证一处登出，处处登出，必须重置 api_token 和 logout_token
// api_token 重置之后会重定向到 sso login page, 但是 sso 依然是登陆状态
// 所以 sso 会重新生成 api_token，导致登出没有效果
// 因此引入 logout_token 如果 sso session 中的 logout_token 不一样，那么代表强制登出
// 以此来保证 用户登出时，api_token 失效，并且 sso 也登陆过期
func (repo *UserRepository) ForceLogout(user *models.User) {
	repo.GenerateApiToken(user, true)
	repo.GenerateLogoutToken(user)
}

func (repo *UserRepository) GenerateApiToken(user *models.User, forceFill bool) string {
	var try int

	log.Debug().Interface("user", user).Msg("GenerateApiToken")
	// 如果生成过token，并且没过期，则直接返回
	if !forceFill &&
		user.ApiToken.Valid &&
		user.ApiToken.String != "" &&
		!repo.TokenExpired(user) {
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

		exists := repo.env.GetDB().Table(models.User{}.TableName()).Where("api_token = ?", str).Find(&models.User{})
		if exists.Error != nil && errors.Is(gorm.ErrRecordNotFound, exists.Error) {
			user.ApiToken = AccessToken
			repo.env.GetDB().Model(user).Update("api_token", AccessToken)
			repo.env.GetDB().Model(user).Update(map[string]interface{}{"api_token": AccessToken, "api_token_created_at": time.Now()})
			return str
		}

		try++
	}
}

func (repo *UserRepository) GenerateLogoutToken(user *models.User) {
	str := helper.RandomString(64)
	repo.env.GetDB().Model(user).Update("logout_token", str)
}

func (repo *UserRepository) GenerateAccessToken(user *models.User) string {
	var (
		try   int
		str   string
		err   error
		reply interface{}
	)

	conn := repo.env.RedisPool().Get()
	defer conn.Close()

	for {
		if try > 10 {
			panic("error GenerateAccessToken try > 10")
		}

		str = helper.RandomString(64)

		reply, err = conn.Do("GET", str)
		if err == nil && reply == nil {
			if repo.env.Config().AccessTokenLifetime > 0 {
				reply, err = conn.Do("SETEX", str, repo.env.Config().AccessTokenLifetime, user.ID)
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

func (repo *UserRepository) FindByToken(token string) (*models.User, error) {
	user := &models.User{}

	if err := repo.env.GetDB().Where("api_token = ?", token).First(user).Error; err != nil {
		log.Debug().Err(err).Msg("FindByToken")
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) TokenExpired(user *models.User) bool {
	seconds := time.Second * time.Duration(repo.env.Config().SessionLifetime)
	if user.ApiToken.Valid &&
		user.ApiToken.String != "" &&
		time.Now().Before(user.ApiTokenCreatedAt.Add(seconds)) {
		return false
	}

	return true
}

func (repo *UserRepository) SyncPermissions(user *models.User, permissions []interface{}) error {
	return repo.env.DBTransaction(func(tx *gorm.DB) error {
		if tx.Model(user).Association("Permissions").Clear().Error != nil {
			return tx.Model(user).Association("Permissions").Clear().Error
		}

		tx.Model(user).Association("Permissions").Append(permissions...)

		return nil
	})
}

func (repo *UserRepository) UpdateLastLoginAt(user *models.User) {
	repo.env.GetDB().Model(user).Update("last_login_at", time.Now())
}

func (repo *UserRepository) FindWithRoles(id int) (*models.User, error) {
	user := &models.User{}

	if err := repo.env.GetDB().Preload("Roles").Where("id = ?", id).First(user).Error; err != nil {
		log.Debug().Err(err).Msg("FindWithRoles")
		return nil, err
	}

	return user, nil
}

func toRoleInterfaceSlice(slice interface{}) []interface{} {
	roles := slice.([]*models.Role)
	newS := make([]interface{}, len(roles))
	for i, v := range roles {
		newS[i] = v
	}

	return newS
}
