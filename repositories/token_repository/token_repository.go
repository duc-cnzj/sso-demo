package token_repository

import "sso/config/env"

type TokenRepositoryImp interface {
}

type TokenRepository struct {
	env *env.Env
}

func NewTokenRepository(env *env.Env) *TokenRepository {
	return &TokenRepository{
		env: env,
	}
}
