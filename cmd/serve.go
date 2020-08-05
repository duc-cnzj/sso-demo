package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
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
		log.Fatal(r.Run(":8888"))
	}()

	<-ctx.Done()
	log.Println("server done by " + ctx.Err().Error())
}

func InitEnv() *env.Env {
	zhLang := zh.New()
	enLang := en.New()
	uni := ut.New(enLang, zhLang, enLang)
	db, err := gorm.Open("mysql", "root:@/sso?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)

	serverEnv := env.NewEnv(db, env.WithUniversalTranslator(uni))

	return serverEnv
}
