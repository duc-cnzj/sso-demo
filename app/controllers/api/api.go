package api

import (
	"github.com/gin-gonic/gin"
	"sso/config/env"
	"sso/repositories/permission_repository"
	"sso/repositories/role_repositories"
	"sso/repositories/token_repository"
	"sso/repositories/user_repository"
)

type AllRepo struct {
	PermRepo *permission_repository.PermissionRepository
	RoleRepo *role_repositories.RoleRepository
	UserRepo *user_repository.UserRepository
	TokenRepo *token_repository.TokenRepository
}

func NewAllRepo(env *env.Env) *AllRepo {
	return &AllRepo{
		PermRepo: permission_repository.NewPermissionRepository(env),
		RoleRepo: role_repositories.NewRoleRepository(env),
		UserRepo: user_repository.NewUserRepository(env),
		TokenRepo: token_repository.NewTokenRepository(env),
	}
}

func Ping(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.String(200, `{"success":true}`)
}

func NotFound(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.String(404, `{"code":404,"message":"Page not found"}`)
}
