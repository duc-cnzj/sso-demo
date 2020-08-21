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
		env := server.Init(".env", "")
		user := models.User{}.FindById(uint(userId), env)
		user.ForceLogout(env)
		log.Println("success")
	}

	//fmt.Println(string(bytes))
}
