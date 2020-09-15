package token_repository

import "sso/config/env"

type TokenRepository struct {
	env *env.Env
}

func NewTokenRepository(env *env.Env) *TokenRepository {
	return &TokenRepository{
		env: env,
	}
}

