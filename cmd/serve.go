package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"sso/server"
	"sso/utils/interrupt"
	"strings"

	"github.com/spf13/cobra"
)

var rootPath string
var ser = &server.Server{}

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

		if err := ser.Init(envPath, rootPath); err != nil {
			log.Fatal().Err(err).Msg("ser.Init")
		}

		go func() {
			log.Fatal().Err(ser.Run()).Msg("server run error")
		}()

		<-ctx.Done()
		ser.Shutdown()
		log.Info().Msg("server done by " + ctx.Err().Error())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&rootPath, "root", ".", "静态资源路径, 必须是dir --root=/path/to/resources")
}
