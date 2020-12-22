package routes

import (
	"net/http"
	"sso/app/controllers/api"
	"sso/app/controllers/api/admin/apitokencontroller"
	adminAuth "sso/app/controllers/api/admin/authcontroller"
	"sso/app/controllers/api/admin/permissioncontroller"
	"sso/app/controllers/api/admin/rolecontroller"
	"sso/app/controllers/api/admin/usercontroller"
	apiWebAuth "sso/app/controllers/api/web/authcontroller"
	webAuth "sso/app/controllers/web/authcontroller"
	webAuthMiddleware "sso/app/middlewares/auth"
	"sso/app/middlewares/i18n"
	"sso/app/middlewares/jwt"
	"sso/config/env"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, env *env.Env) *gin.Engine {
	router.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders(webAuthMiddleware.HttpAuthToken, "Authorization")
	router.Use(cors.New(config))
	router.Use(sessions.Sessions("sso", env.SessionStore()), i18n.I18nMiddleware(env))

	if !env.IsSkipLoadResources() {
		router.Static("/assets", env.RootDir()+"resources/views/web/css")
		router.Static("/images", env.RootDir()+"resources/views/web/images")
		router.Static("/avatars", env.RootDir()+"resources/views/admin/avatars")

		router.Static("/static/css", env.RootDir()+"resources/views/admin/static/css")
		router.Static("/static/fonts", env.RootDir()+"resources/views/admin/static/fonts")
		router.Static("/static/img", env.RootDir()+"resources/views/admin/static/img")
		router.Static("/static/js", env.RootDir()+"resources/views/admin/static/js")
		router.LoadHTMLFiles(env.RootDir()+"resources/views/web/login.tmpl", env.RootDir()+"resources/views/web/select_system.tmpl", env.RootDir()+"resources/views/admin/index.html")
		router.StaticFile("/favicon.ico", env.RootDir()+"resources/views/admin/favicon.ico")
	}

	// for debug
	//router.LoadHTMLGlob("/Users/congcong/uco/sso/resources/views/*")

	router.NoRoute(api.NotFound)

	router.Any("/ping", api.Ping)

	auth := webAuth.New(env)

	guest := router.Group("/", webAuthMiddleware.GuestMiddleware(env))
	{
		guest.GET("/login", auth.LoginForm)
		guest.POST("/login", auth.Login)
	}

	authRouter := router.Group("/", webAuthMiddleware.SessionMiddleware(env))
	{
		authRouter.GET("/", auth.SelectSystem)
		authRouter.POST("/auth/logout", auth.Logout)
	}

	router.POST("/access_token", auth.AccessToken)

	webApiGroup := router.Group("/api", webAuthMiddleware.ApiMiddleware(env))
	{
		webApiAuth := apiWebAuth.NewAuthController(env)
		webApiGroup.POST("/user/info", webApiAuth.Info)
		webApiGroup.POST("/user/info/projects/:project", webApiAuth.Info)
		webApiGroup.POST("/logout", webApiAuth.Logout)
	}

	adminGroup := router.Group("/api/admin")
	{
		apiAuth := adminAuth.New(env)
		adminGroup.POST("/login", apiAuth.Login)

		apiGroup := adminGroup.Group("/", jwt.AuthMiddleware(env))
		apiGroup.POST("/user/info", apiAuth.Info)
		apiGroup.POST("/logout", apiAuth.Logout)

		role := rolecontroller.NewRoleController(env)
		apiGroup.GET("/all_roles", role.All)
		apiGroup.GET("/roles", role.Index)
		apiGroup.POST("/roles", role.Store)
		apiGroup.GET("/roles/:role", role.Show)
		apiGroup.PUT("/roles/:role", role.Update)
		apiGroup.DELETE("/roles/:role", role.Destroy)

		permissions := permissioncontroller.NewPermissionController(env)
		apiGroup.GET("/permissions_by_group", permissions.GetByGroups)
		apiGroup.GET("/get_permission_projects", permissions.GetPermissionProjects)
		apiGroup.GET("/permissions", permissions.Index)
		apiGroup.POST("/permissions", permissions.Store)
		apiGroup.GET("/permissions/:permission", permissions.Show)
		apiGroup.PUT("/permissions/:permission", permissions.Update)
		apiGroup.DELETE("/permissions/:permission", permissions.Destroy)

		user := usercontroller.NewUserController(env)
		apiGroup.POST("/users/:user/change_password", user.ChangePassword)
		apiGroup.POST("/users/:user/force_logout", user.ForceLogout)
		apiGroup.POST("/users/:user/sync_roles", user.SyncRoles)
		apiGroup.GET("/users", user.Index)
		apiGroup.POST("/users", user.Store)
		apiGroup.GET("/users/:user", user.Show)
		apiGroup.PUT("/users/:user", user.Update)
		apiGroup.DELETE("/users/:user", user.Destroy)

		apiToken := apitokencontroller.New(env)
		apiGroup.GET("/users/:user/api_tokens", apiToken.Index)
		apiGroup.GET("/api_tokens", apiToken.Index)
	}

	router.Any("/admin/*action", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	return router
}
