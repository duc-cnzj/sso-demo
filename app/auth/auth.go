package auth

import (
	"sso/app/models"
)

type Auth struct {
	user *models.User
}

func (auth *Auth) SetUser(user *models.User) {
	auth.user = user
}

func NewAuth(user *models.User) *Auth {
	return &Auth{user: user}
}

func (auth *Auth) HasPermission(name string) bool {
	if auth.user == nil {
		return false
	}
	for _, perm := range auth.user.CurrentPermissions {
		if perm == name {
			return true
		}
	}

	return false
}

func (auth *Auth) HasRole(name string) bool {
	if auth.user == nil {
		return false
	}
	for _, role := range auth.user.CurrentRoles {
		if role == name {
			return true
		}
	}

	return false
}

func (auth *Auth) IsAdmin() bool {
	var adminEmails = []string{"1025434218@qq.com"}

	if auth.user == nil {
		return false
	}

	if auth.user.ID == 1 {
		return true
	}

	for _, email := range adminEmails {
		if auth.user.Email == email {
			return true
		}
	}

	return false
}
