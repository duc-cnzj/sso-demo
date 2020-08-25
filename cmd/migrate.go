package cmd

import (
	"github.com/jinzhu/gorm"
	"log"
	"sso/app/models"
	"sso/config/env"
	"sso/server"

	"github.com/spf13/cobra"
)

var migrateModels = []interface{}{
	&models.User{},
	&models.Role{},
	&models.Permission{},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "数据库迁移",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err    error
			config env.Config
			conn   *gorm.DB
		)
		if config, err = server.ReadConfig(envPath); err != nil {
			log.Panicln(err)
		}
		conn, err = server.DB(config)
		if err != nil {
			log.Fatal(err)
		}
		migrate := conn.AutoMigrate(migrateModels...)
		if migrate.Error != nil {
			log.Fatal("migrate.Error", migrate.Error.Error())
		}

		log.Println("migrate ok!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
