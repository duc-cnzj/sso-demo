package tests

import (
	"github.com/gin-gonic/gin"
	"sso/config/env"
	"sso/routes"
	"sso/server"
)

func InitRouter(env *env.Env) *gin.Engine {
	r := gin.New()

	return routes.Init(r, env)
}

func NewTestEnv(path string) *env.Env {
	return server.Init(path, "")
}
