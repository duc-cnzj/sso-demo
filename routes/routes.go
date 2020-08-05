package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"sso/app/controllers"
	auth2 "sso/app/http/middlewares/auth"
	"sso/app/http/middlewares/i18n"
	"sso/config/env"
)

func Init(router *gin.Engine, env *env.Env)  {
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	router.Use(sessions.Sessions("sso", store), i18n.I18nMiddleware(env))

	router.LoadHTMLGlob("resources/views/*")

	auth := authcontroller.New(env)

	guest := router.Group("/", auth2.GuestMiddleware(env))
	{
		guest.GET("/login", auth.LoginForm)
		guest.POST("/login", auth.Login)
	}

	// get /login?redirect_url=xxxxx
	// 返回 code

	// post /oauth/access_token
	// 返回 token

	//authRouter := router.Use(auth2.AuthMiddleware(env))
	authRouter := router.Group("/auth", auth2.SessionMiddleware(env))
	{
		// get /auth/user 带上token
		// 返回用户信息
		authRouter.GET("/user", auth.Me)
		authRouter.GET("/logout", auth.Logout)
	}
}