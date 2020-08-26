package server

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	_ "github.com/go-sql-driver/mysql"
	redis2 "github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"path"
	"sso/app/models"
	"sso/config/env"
	"time"
)

func Init(configPath string, rootPath string) *env.Env {
	var (
		config env.Config
		err    error
		db     *gorm.DB
	)

	if config, err = ReadConfig(configPath); err != nil {
		return nil
	}

	if config.Debug {
		fmt.Printf(`
			AppPort:             %d,
			Debug:               %t,
			DBConnection:        %s,
			DBHost:              %s,
			DBPort:              %d,
			DBDatabase:          %s,
			DBUsername:          %s,
			DBPassword:          %s,
			RedisHost:           %s,
			RedisPassword:       %s,
			RedisPort:           %d,
			SessionLifetime:     %d,
			AccessTokenLifetime: %d,
			JwtSecret: 		     %s,
			JwtExpiresSeconds:        %d,
`,
			config.AppPort,
			config.Debug,
			config.DBConnection,
			config.DBHost,
			config.DBPort,
			config.DBDatabase,
			config.DBUsername,
			config.DBPassword,
			config.RedisHost,
			config.RedisPassword,
			config.RedisPort,
			config.SessionLifetime,
			config.AccessTokenLifetime,
			config.JwtSecret,
			config.JwtExpiresSeconds,
		)
	}

	zhLang := zh.New()
	enLang := en.New()
	uni := ut.New(enLang, zhLang, enLang)

	if db, err = DB(config); err != nil {
		log.Panicln(err)
	}

	password := config.RedisPassword
	redisPool := redisPool(config, password)

	store := sessionStore(redisPool, config)
	serverEnv := env.NewEnv(config, db, store, redisPool, env.WithUniversalTranslator(uni), env.WithRootDir(rootPath))
	gob.Register(&models.User{})

	return serverEnv
}

func ReadConfig(configPath string) (env.Config, error) {
	var err error
	if configPath == "" {
		configPath = ".env"
	}
	if !path.IsAbs(configPath) {
		viper.AddConfigPath(".")
	}
	viper.SetConfigFile(configPath)

	if err = viper.ReadInConfig(); err != nil {
		return env.Config{}, err
	}

	config := env.Config{
		AppPort:             viper.GetUint("APP_PORT"),
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
		JwtSecret:           viper.GetString("JWT_SECRET"),
		JwtExpiresSeconds:   viper.GetInt64("JWT_EXPIRES_AT"),
	}

	return config, nil
}

func DB(config env.Config) (*gorm.DB, error) {
	db, err := gorm.Open(config.DBConnection, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBDatabase))
	if err != nil {
		return db, err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)
	if config.Debug {
		db.LogMode(true)
	}
	return db, nil
}

func sessionStore(redisPool *redis2.Pool, config env.Config) redis.Store {
	store, err := redis.NewStoreWithPool(redisPool, []byte("secret"))
	if err != nil {
		log.Panicln(err)
		return nil
	}
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: config.SessionLifetime,
	})
	return store
}

func redisPool(config env.Config, password string) *redis2.Pool {
	redisPool := &redis2.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis2.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis2.Conn, error) {
			c, err := redis2.Dial("tcp", fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort))
			if err != nil {
				log.Panicln(err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", config.RedisPassword); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
	}
	return redisPool
}
