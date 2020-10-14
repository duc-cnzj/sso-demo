package cmd

import (
	"log"
	"sso/app/models"
	"sso/server"

	"github.com/spf13/cobra"
)

var migrateModels = []interface{}{
	&models.User{},
	&models.Role{},
	&models.Permission{},
	&models.ApiToken{},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "数据库迁移",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			s   = &server.Server{}
		)
		s.SetRunningInConsole()
		if err = s.Init(envPath, ""); err != nil {
			return
		}
		migrate := s.Env().GetDB().AutoMigrate(migrateModels...)
		if migrate.Error != nil {
			log.Fatal("migrate.Error", migrate.Error.Error())
		}

		log.Println("migrate ok!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
