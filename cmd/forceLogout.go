package cmd

import (
	"log"
	"sso/repositories/user_repository"
	"sso/server"

	"github.com/spf13/cobra"
)

var userId uint

var forceLogoutCmd = &cobra.Command{
	Use:   "forceLogout",
	Short: "强制用户登出",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			s   = &server.Server{}
		)

		if err = s.Init(envPath, ""); err != nil {
			return
		}

		if userId > 0 {
			var env = s.Env()
			userRepo := user_repository.NewUserRepository(env)
			user, _ := userRepo.FindById(userId)
			userRepo.ForceLogout(user)
			log.Println("success")
		}
	},
}

func init() {
	rootCmd.AddCommand(forceLogoutCmd)

	forceLogoutCmd.Flags().UintVarP(&userId, "userId", "u", 0, "-u|--userId=1")
}
