package api

import (
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
