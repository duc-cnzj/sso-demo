package server

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	_ "github.com/go-sql-driver/mysql"
	redis2 "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"sso/app/models"
	"sso/config/env"
	"sso/routes"
	"time"
)

type Server struct {
	env     *env.Env
	config  *env.Config
	db      *gorm.DB
	redis   *redis2.Pool
	session sessions.Store
	engine  *gin.Engine
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) Config() *env.Config {
	return s.config
}

func (s *Server) Env() *env.Env {
	return s.env
}

func (s *Server) Init(configPath, rootPath string) error {
	if err := s.LoadConfig(configPath); err != nil {
		return err
	}

	if err := s.InitDB(); err != nil {
		return err
	}

	s.InitRedis()

	if err := s.InitSession(); err != nil {
		return err
	}

	zhLang := zh.New()
	enLang := en.New()
	uni := ut.New(enLang, zhLang, enLang)

	s.env = env.NewEnv(
		s.config,
		s.db,
		s.session,
		s.redis,
		env.WithUniversalTranslator(uni),
		env.WithRootDir(rootPath),
	)

	gob.Register(&models.User{})

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if s.env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if s.env.IsDebugging() {
		gin.SetMode(gin.DebugMode)
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			log.Debug().Msgf("route:%10s\t%v", httpMethod, absolutePath)
		}

		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("############### debug mode ###############")
		s.env.PrintConfig()
		s.db.LogMode(true)
	}

	r := gin.New()

	s.engine = routes.Init(r, s.env)

	return nil
}

func (s *Server) Run() error {
	return s.engine.Run(fmt.Sprintf(":%d", s.config.AppPort))
}

func (s *Server) Shutdown() {

}

func (s *Server) LoadConfig(configPath string) error {
	var (
		config *env.Config
		err    error
	)
	if config, err = ReadConfig(configPath); err != nil {
		return err
	}
	s.config = config

	return nil
}

func (s *Server) InitDB() error {
	var err error
	s.db, err = gorm.Open(s.config.DBConnection, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci", s.config.DBUsername, s.config.DBPassword, s.config.DBHost, s.config.DBPort, s.config.DBDatabase))
	if err != nil {
		return err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	s.db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	s.db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	s.db.DB().SetConnMaxLifetime(time.Hour)

	return nil
}

func (s *Server) InitRedis() {
	s.redis = &redis2.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis2.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis2.Conn, error) {
			c, err := redis2.Dial("tcp", fmt.Sprintf("%s:%d", s.config.RedisHost, s.config.RedisPort))

			if err != nil {
				return nil, err
			}

			if s.config.RedisPassword != "" {
				if _, err := c.Do("AUTH", s.config.RedisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
	}
}

func (s *Server) InitSession() error {
	var (
		store redis.Store
		err   error
	)
	if store, err = redis.NewStoreWithPool(s.redis, []byte("secret")); err != nil {
		return err
	}

	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: s.config.SessionLifetime,
	})
	s.session = store

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadConfig(configPath string) (*env.Config, error) {
	var err error
	if configPath == "" {
		configPath = ".env"
	}
	if !path.IsAbs(configPath) {
		getwd, _ := os.Getwd()
		configPath = getwd + "/" + configPath
	}
	exists := fileExists(configPath)
	if !exists {
		return nil, errors.New("file not exists in " + configPath)
	}
	viper.SetConfigType("env")
	file, _ := ioutil.ReadFile(configPath)

	if err = viper.ReadConfig(bytes.NewReader(file)); err != nil {
		return &env.Config{}, err
	}

	config := &env.Config{
		AppPort:             viper.GetUint("APP_PORT"),
		AppEnv:              viper.GetString("APP_ENV"),
		Debug:               viper.GetBool("DEBUG"),
		DBConnection:        viper.GetString("DB_CONNECTION"),
		DBHost:              viper.GetString("DB_HOST"),
		DBPort:              viper.GetUint("DB_PORT"),
		DBDatabase:          viper.GetString("DB_DATABASE"),
		DBUsername:          viper.GetString("DB_USERNAME"),
		DBPassword:          viper.GetString("DB_PASSWORD"),
		RedisHost:           viper.GetString("REDIS_HOST"),
		RedisPassword:       viper.GetString("REDIS_PASSWORD"),
		RedisPort:           viper.GetUint("REDIS_PORT"),
		SessionLifetime:     viper.GetInt("SESSION_LIFETIME"),
		AccessTokenLifetime: viper.GetInt("ACCESS_TOKEN_LIFETIME"),
		ApiTokenLifetime:    viper.GetInt("API_TOKEN_LIFETIME"),
		JwtSecret:           viper.GetString("JWT_SECRET"),
		JwtExpiresSeconds:   viper.GetInt64("JWT_EXPIRES_AT"),
	}

	return config, nil
}
