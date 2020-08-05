package routes

import (
	"github.com/gin-gonic/gin"
	"sso/app/controllers"
	auth2 "sso/app/http/middlewares/auth"
	"sso/app/http/middlewares/i18n"
	"sso/config/env"
)

func Init(router *gin.Engine, env *env.Env)  {
	router.Use(i18n.I18nMiddleware(env))

	router.LoadHTMLGlob("resources/views/*")

	auth := authcontroller.New(env)

	router.GET("/login", auth.LoginForm)

	router.POST("/login", auth.Login)

	// get /login?redirect_url=xxxxx
	// 返回 code

	// post /oauth/access_token
	// 返回 token

	authRouter := router.Use(auth2.AuthMiddleware(env))
	// get /auth/user 带上token
	// 返回用户信息
	authRouter.GET("/auth/user", auth.Me)
}