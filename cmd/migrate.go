package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sso/app/models"
)

var migrateModels = []interface{}{
	&models.User{},
	&models.Role{},
	&models.Permission{},
}

func main() {
	db, err := gorm.Open("mysql", "root:@/sso?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	migrate := db.AutoMigrate(migrateModels...)
	if migrate.Error != nil {
		log.Fatal(migrate.Error.Error())
	}

	log.Println("migrate ok!")
}
