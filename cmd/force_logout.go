package main

import (
	"flag"
	"log"
	"sso/app/models"
	"sso/server"
)

var userId int

func init() {
	flag.IntVar(&userId, "id", 0, "-id")
}

func main() {
	flag.Parse()

	if userId > 0 {
		env := server.Init()
		user := models.User{}.FindById(uint(userId), env)
		user.GenerateLogoutToken(env)
		user.GenerateApiToken(env, true)
		log.Println("success")
	}

	//fmt.Println(string(bytes))
}
