package routes

import (
	"github.com/gin-gonic/gin"
	"sso/app/controllers"
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

	// get /auth/user 带上token
	// 返回用户信息
}