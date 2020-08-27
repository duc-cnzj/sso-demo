package cmd

import (
	"github.com/jinzhu/gorm"
	"log"
	"sso/config/env"
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
			err    error
			config env.Config
			db     *gorm.DB
		)

		if config, err = server.ReadConfig(envPath); err != nil {
			return
		}

		if db, err = server.DB(config); err != nil {
			return
		}

		if userId > 0 {
			var env = env.NewEnv(config, nil, nil, nil, env.WithDB(db))
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
