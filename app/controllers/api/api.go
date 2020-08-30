package api

import (
	"github.com/gin-gonic/gin"
	"sso/config/env"
	"sso/repositories/permission_repository"
	"sso/repositories/role_repositories"
	"sso/repositories/user_repository"
)

type AllRepo struct {
	PermRepo *permission_repository.PermissionRepository
	RoleRepo *role_repositories.RoleRepository
	UserRepo *user_repository.UserRepository
}

func NewAllRepo(env *env.Env) *AllRepo {
	return &AllRepo{
		PermRepo: permission_repository.NewPermissionRepository(env),
		RoleRepo: role_repositories.NewRoleRepository(env),
		UserRepo: user_repository.NewUserRepository(env),
	}
}

func Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"success": true})
}

func NotFound(c *gin.Context) {
	c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
}
