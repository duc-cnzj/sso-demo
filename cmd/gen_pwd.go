package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var password string

func init() {
	flag.StringVar(&password, "password", "", "-password")
}

func main() {
	flag.Parse()
	bytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println(string(bytes))
}
