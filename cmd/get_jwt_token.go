package main

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"sso/app/models"
	"sso/config/env"
	_ "github.com/go-sql-driver/mysql"
	"sso/utils/helper"
	"sso/utils/jwt"
)

var id uint

func init() {
	flag.UintVar(&id, "id", 0, "-id")
}

func main()  {
	flag.Parse()
	db, e := gorm.Open("mysql", "root:@/sso?charset=utf8mb4&parseTime=True&loc=Local")
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	fmt.Println(jwt.GenerateToken(models.User{}.FindById(id, env.NewEnv(db))))
}