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
	DBConnection string
	DBHost       string
	DBPort       uint
	DBDatabase   string
	DBUsername   string
	DBPassword   string

	// redis
	RedisHost     string
	RedisPassword string
	RedisPort     uint
	Debug         bool

	// jwt
	JwtSecret         string
	JwtExpiresSeconds int64
}

type Env struct {
	translator   *ut.UniversalTranslator
	db           *gorm.DB
	sessionStore sessions.Store
	redisPool    *redis.Pool
	config       Config
	rootDir      string
}

func (e *Env) IsDebugging() bool {
	return e.config.Debug
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

func WithDB(db *gorm.DB) func(env *Env) {
	return func(env *Env) {
		env.db = db
	}
}

func WithRootDir(path string) func(env *Env) {
	return func(env *Env) {
		env.rootDir = path
	}
}

func (e *Env) GetUniversalTranslator() *ut.UniversalTranslator {
	return e.translator
}

func (e *Env) GetDB() *gorm.DB {
	return e.db
}

func (e *Env) DBTransaction(fn func(tx *gorm.DB) error) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (e *Env) RootDir() string {
	return e.rootDir
}
