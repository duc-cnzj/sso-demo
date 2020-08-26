package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sso/routes"
	"sso/server"
	"sso/utils/interrupt"
	"strings"

	"github.com/spf13/cobra"
)

var rootPath string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 sso 服务。",
	PreRun: func(cmd *cobra.Command, args []string) {
		if strings.HasSuffix(rootPath, "/") {
			rootPath = strings.TrimSuffix(rootPath, "/")
		}
		rootPath = rootPath + "/"

		fmt.Println("env file path: ", envPath)
		fmt.Println("root path: ", rootPath)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx, done := interrupt.Context()

		defer done()

		serverEnv := server.Init(envPath, rootPath)

		if !serverEnv.IsDebugging() {
			gin.SetMode(gin.ReleaseMode)
		}

		r := gin.Default()

		routes.Init(r, serverEnv)
		go func() {
			log.Fatal(r.Run(fmt.Sprintf(":%d", serverEnv.Config().AppPort)))
		}()

		<-ctx.Done()
		log.Println("server done by " + ctx.Err().Error())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&rootPath, "root", ".", "静态资源路径, 必须是dir --root=/path/to/resources")
}
