package env

import (
	"github.com/gin-contrib/sessions"
	ut "github.com/go-playground/universal-translator"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type Config struct {
	AppPort             uint
	SessionLifetime     int
	AccessTokenLifetime int

	// db
	DbConnection string
	DbHost       string
	DbPort       uint
	DbDatabase   string
	DbUsername   string
	DbPassword   string

	// redis
	RedisHost     string
	RedisPassword string
	RedisPort     uint
}

type Env struct {
	translator   *ut.UniversalTranslator
	db           *gorm.DB
	sessionStore sessions.Store
	redisPool    *redis.Pool
	config       Config
}

func (e *Env) Config() Config {
	return e.config
}

func (e *Env) RedisPool() *redis.Pool {
	return e.redisPool
}

func (e *Env) SessionStore() sessions.Store {
	return e.sessionStore
}

type Operator func(*Env)

func NewEnv(config Config, db *gorm.DB, sessionStore sessions.Store, pool *redis.Pool, envOperators ...Operator) *Env {
	env := &Env{
		db:           db,
		sessionStore: sessionStore,
		redisPool:    pool,
		config:       config,
	}
	for _, op := range envOperators {
		op(env)
	}

	return env
}

func WithUniversalTranslator(t *ut.UniversalTranslator) func(env *Env) {
	return func(env *Env) {
		env.translator = t
	}
}

func (e *Env) GetUniversalTranslator() *ut.UniversalTranslator {
	return e.translator
}

func (e *Env) GetDB() *gorm.DB {
	return e.db
}
