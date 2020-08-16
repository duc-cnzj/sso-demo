package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sso/app/http/controllers/authcontroller"
	"sso/app/http/controllers/permissioncontroller"
	"sso/app/http/controllers/rolecontroller"
	auth2 "sso/app/http/middlewares/auth"
	"sso/app/http/middlewares/i18n"
	"sso/config/env"
)

func Init(router *gin.Engine, env *env.Env) {
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

	api := router.Group("/api")
	//api := router.Group("/api", auth2.ApiMiddleware(env))
	{
		api.POST("/user/info", auth.Info)

		role := rolecontroller.NewRoleController(env)
		api.GET("/roles", role.Index)
		api.POST("/roles", role.Store)
		api.GET("/roles/:role", role.Show)
		api.PUT("/roles/:role", role.Update)
		api.DELETE("/roles/:role", role.Destroy)

		permissions := permissioncontroller.NewPermissionController(env)
		api.GET("/permissions", permissions.Index)
		api.POST("/permissions", permissions.Store)
		api.GET("/permissions/:permission", permissions.Show)
		api.PUT("/permissions/:permission", permissions.Update)
		api.DELETE("/permissions/:permission", permissions.Destroy)
	}

}
