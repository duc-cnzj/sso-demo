package jwt

import (
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog/log"
	"sso/config/env"
	"time"
)

const KeyPrefix = "jwt_blacklist:"

func AddToBlacklist(seconds int64, token string, env *env.Env) {
	conn := env.RedisPool().Get()
	defer conn.Close()
	k := getKey(token)
	do, err := conn.Do("SETEX", k, seconds, time.Now().Unix())
	log.Debug().Err(err).Interface("do", do).Msg("AddToBlacklist")
}

func KeyInBlacklist(token string, env *env.Env) bool {
	log.Debug().Msg("KeyInBlacklist: " + token)
	conn := env.RedisPool().Get()
	defer conn.Close()
	k := getKey(token)
	do, err := redis.String(conn.Do("GET", k))
	log.Error().Err(err).Msg("jwt.KeyInBlacklist")
	if do != "" {
		return true
	}

	return false
}

func getKey(token string) string {
	k := KeyPrefix + token
	return k
}
