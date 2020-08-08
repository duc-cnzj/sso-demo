package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sso/routes"
	"sso/server"
	"sso/utils/interrupt"
)

func main() {
	ctx, done := interrupt.Context()
	defer done()

	r := gin.Default()

	serverEnv := server.Init()

	routes.Init(r, serverEnv)
	go func() {
		log.Fatal(r.Run(fmt.Sprintf(":%d", serverEnv.Config().AppPort)))
	}()

	<-ctx.Done()
	log.Println("server done by " + ctx.Err().Error())
}
