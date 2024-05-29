package global

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

var (
	HOST string = "localhost"
	PORT string = "8080"
)

func SetEnv() {
	// Env variables
	// get parameter from .env file
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	if myEnv["HOST"] != "" {
		HOST = myEnv["HOST"]
	}
	if myEnv["PORT"] != "" {
		PORT = myEnv["PORT"]
	}
	fmt.Println(HOST, PORT)
}
