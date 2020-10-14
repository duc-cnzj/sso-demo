package cmd

import (
	"fmt"
	"sso/app/middlewares/jwt"
	"sso/server"

	"github.com/spf13/cobra"
)

var jwtString string
var parseJwtCmd = &cobra.Command{
	Use:   "parseJwt",
	Short: "解析jwt",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			s   = &server.Server{}
		)
		s.SetRunningInConsole()
		if err = s.Init(envPath, ""); err != nil {
			return
		}
		parse, e := jwt.Parse(jwtString, s.Env())
		fmt.Println(parse, e)
	},
}

func init() {
	rootCmd.AddCommand(parseJwtCmd)

	parseJwtCmd.Flags().StringVar(&jwtString, "jwt", "", "--jwt")
}
