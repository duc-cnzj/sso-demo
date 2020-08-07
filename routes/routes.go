package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sso/app/controllers"
	auth2 "sso/app/http/middlewares/auth"
	"sso/app/http/middlewares/i18n"
	"sso/config/env"
)

func Init(router *gin.Engine, env *env.Env)  {
	router.Use(sessions.Sessions("sso", env.SessionStore()), i18n.I18nMiddleware(env))

	router.Static("/assets", "resources/css")
	router.Static("/images", "resources/images")
	router.LoadHTMLGlob("resources/views/*")
	// for debug
	//router.LoadHTMLGlob("/Users/congcong/uco/sso/resources/views/*")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Page not found"})
	})

	auth := authcontroller.New(env)

	guest := router.Group("/", auth2.GuestMiddleware(env))
	{
		guest.GET("/login", auth.LoginForm)
		guest.POST("/login", auth.Login)
	}

	authRouter := router.Group("/auth", auth2.SessionMiddleware(env))
	{
		authRouter.GET("/select_system", auth.SelectSystem)
		authRouter.GET("/logout", auth.Logout)
	}

	router.POST("/access_token", auth.AccessToken)
	router.POST("/user/info", auth.Info)
}