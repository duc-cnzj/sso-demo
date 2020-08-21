package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sso/routes"
	"sso/server"
	"sso/utils/interrupt"
)

var envPath string
var rootPath string

func init() {
	flag.StringVar(&envPath, "config", ".env", "-config")
	flag.StringVar(&rootPath, "root", "", "-root")
}

func main() {
	flag.Parse()
	ctx, done := interrupt.Context()
	defer done()

	r := gin.Default()

	serverEnv := server.Init(envPath, rootPath)

	routes.Init(r, serverEnv)
	go func() {
		log.Fatal(r.Run(fmt.Sprintf(":%d", serverEnv.Config().AppPort)))
	}()

	<-ctx.Done()
	log.Println("server done by " + ctx.Err().Error())
}
