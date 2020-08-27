package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
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

		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		serverEnv := server.Init(envPath, rootPath)

		gin.SetMode(gin.ReleaseMode)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		if !serverEnv.IsProduction() {
			gin.SetMode(gin.DebugMode)
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
			log.Info().Msg("############### debug mode ###############")
		}

		r := gin.New()
		//r := gin.Default()

		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
			log.Debug().Msgf("route:%10s\t%v", httpMethod, absolutePath)
		}

		routes.Init(r, serverEnv)

		go func() {
			log.Fatal().Err(r.Run(fmt.Sprintf(":%d", serverEnv.Config().AppPort))).Msg("server run error")
		}()

		<-ctx.Done()
		log.Info().Msg("server done by " + ctx.Err().Error())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&rootPath, "root", ".", "静态资源路径, 必须是dir --root=/path/to/resources")
}
