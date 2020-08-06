package main

import (
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
	"github.com/spf13/viper"
	"time"

	"log"
	"sso/config/env"
	"sso/routes"
	"sso/utils/interrupt"
)

func main() {
	ctx, done := interrupt.Context()
	defer done()

	r := gin.Default()

	serverEnv := InitEnv()

	routes.Init(r, serverEnv)
	go func() {
		log.Fatal(r.Run(fmt.Sprintf(":%d", serverEnv.Config().AppPort)))
	}()

	<-ctx.Done()
	log.Println("server done by " + ctx.Err().Error())
}

func InitEnv() *env.Env {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env") // REQUIRED if the config file does not have the extension in the name
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	config := env.Config{
		AppPort:         viper.GetUint("APP_PORT"),
		DbConnection:    viper.GetString("DB_CONNECTION"),
		DbHost:          viper.GetString("DB_HOST"),
		DbPort:          viper.GetUint("DB_PORT"),
		DbDatabase:      viper.GetString("DB_DATABASE"),
		DbUsername:      viper.GetString("DB_USERNAME"),
		DbPassword:      viper.GetString("DB_PASSWORD"),
		RedisHost:       viper.GetString("REDIS_HOST"),
		RedisPassword:   viper.GetString("REDIS_PASSWORD"),
		RedisPort:       viper.GetUint("REDIS_PORT"),
		SessionLifetime: time.Duration(viper.GetInt64("SESSION_LIFETIME")),
		AccessTokenLifetime: time.Duration(viper.GetInt64("ACCESS_TOKEN_LIFETIME")),
	}
	fmt.Println(config)

	zhLang := zh.New()
	enLang := en.New()
	uni := ut.New(enLang, zhLang, enLang)
	db, err := gorm.Open(config.DbConnection, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbDatabase))
	if err != nil {
		panic(err)
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)

	//rediStore, err := redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	password := ""

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
	store, _ := redis.NewStoreWithPool(redisPool, []byte("secret"))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: int(config.SessionLifetime),
	})
	serverEnv := env.NewEnv(config, db, store, redisPool, env.WithUniversalTranslator(uni))

	return serverEnv
}
