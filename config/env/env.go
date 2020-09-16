package env

import (
	"github.com/gin-contrib/sessions"
	ut "github.com/go-playground/universal-translator"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"sync"
)

type Config struct {
	AppEnv              string
	AppPort             uint
	SessionLifetime     int
	ApiTokenLifetime    int
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
	config       *Config
	rootDir      string
	mu           *sync.Mutex
}

func (e *Env) SetDB(db *gorm.DB) {
	e.db = db
}

func (e *Env) IsDebugging() bool {
	return e.config.Debug
}

func (e *Env) Config() *Config {
	return e.config
}

func (e *Env) RedisPool() *redis.Pool {
	return e.redisPool
}

func (e *Env) IsProduction() bool {
	return e.Config().AppEnv == "production" || e.Config().AppEnv == "prod" || e.Config().AppEnv == ""
}

func (e *Env) IsLocal() bool {
	return e.Config().AppEnv == "local"
}

func (e *Env) IsTesting() bool {
	return e.Config().AppEnv == "testing"
}

func (e *Env) SessionStore() sessions.Store {
	return e.sessionStore
}

type Operator func(*Env)

func NewEnv(config *Config, db *gorm.DB, sessionStore sessions.Store, pool *redis.Pool, envOperators ...Operator) *Env {
	env := &Env{
		db:           db,
		sessionStore: sessionStore,
		redisPool:    pool,
		config:       config,
		mu:           &sync.Mutex{},
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

func (e *Env) PrintConfig() {
	log.Debug().Msgf("%20s: %v", "AppPort", e.config.AppPort)
	log.Debug().Msgf("%20s: %v", "AppEnv", e.config.AppEnv)
	log.Debug().Msgf("%20s: %v", "Debug", e.config.Debug)
	log.Debug().Msgf("%20s: %v", "DBConnection", e.config.DBConnection)
	log.Debug().Msgf("%20s: %v", "DBHost", e.config.DBHost)
	log.Debug().Msgf("%20s: %v", "DBPort", e.config.DBPort)
	log.Debug().Msgf("%20s: %v", "DBDatabase", e.config.DBDatabase)
	log.Debug().Msgf("%20s: %v", "DBUsername", e.config.DBUsername)
	log.Debug().Msgf("%20s: %v", "DBPassword", e.config.DBPassword)
	log.Debug().Msgf("%20s: %v", "RedisHost", e.config.RedisHost)
	log.Debug().Msgf("%20s: %v", "RedisPassword", e.config.RedisPassword)
	log.Debug().Msgf("%20s: %v", "RedisPort", e.config.RedisPort)
	log.Debug().Msgf("%20s: %v", "SessionLifetime", e.config.SessionLifetime)
	log.Debug().Msgf("%20s: %v", "AccessTokenLifetime", e.config.AccessTokenLifetime)
	log.Debug().Msgf("%20s: %v", "ApiTokenLifetime", e.config.ApiTokenLifetime)
	log.Debug().Msgf("%20s: %v", "JwtSecret", e.config.JwtSecret)
	log.Debug().Msgf("%20s: %v", "JwtExpiresSeconds", e.config.JwtExpiresSeconds)
}
