package cmd

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"

	"github.com/spf13/cobra"
)

var password string

var genPwdCmd = &cobra.Command{
	Use:   "genPwd",
	Short: "生成密码",
	Run: func(cmd *cobra.Command, args []string) {
		bytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			log.Fatal(e)
		}

		fmt.Println(string(bytes))
	},
}

func init() {
	rootCmd.AddCommand(genPwdCmd)
	genPwdCmd.Flags().StringVarP(&password, "password", "p", "", "-p|--password")
}
