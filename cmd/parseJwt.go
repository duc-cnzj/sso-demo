/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"sso/app/http/middlewares/jwt"
	"sso/config/env"
	"sso/server"

	"github.com/spf13/cobra"
)

var jwtString string
var parseJwtCmd = &cobra.Command{
	Use:   "parseJwt",
	Short: "解析jwt",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err    error
			config env.Config
		)
		if config, err = server.ReadConfig(envPath); err != nil {
			return
		}
		env := env.NewEnv(config, nil, nil, nil)
		parse, e := jwt.Parse(jwtString, env)
		fmt.Println(parse, e)
	},
}

func init() {
	rootCmd.AddCommand(parseJwtCmd)

	parseJwtCmd.Flags().StringVar(&jwtString, "jwt", "", "--jwt")
}
