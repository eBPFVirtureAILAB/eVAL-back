package global

import (
	"log"

	"github.com/joho/godotenv"
)

var (
	HOST     string = "localhost"
	PORT     string = "8080"
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
)

func SetEnv() {
	// initialize variables
	Postgres.Host = "db"
	Postgres.Port = "5432"
	Postgres.User = "postgres"
	Postgres.Password = "postgres"
	Postgres.DBName = "eval"

	// Env variables
	// get parameter from .env file
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Printf("Error loading .env file")
		return
	}

	// Web server
	if myEnv["HOST"] != "" {
		HOST = myEnv["HOST"]
	}
	if myEnv["PORT"] != "" {
		PORT = myEnv["PORT"]
	}

	// Postgres
	if myEnv["POSTGRES_HOST"] != "" {
		Postgres.Host = myEnv["POSTGRES_HOST"]
	}
	if myEnv["POSTGRES_PORT"] != "" {
		Postgres.Port = myEnv["POSTGRES_PORT"]
	}
	if myEnv["POSTGRES_USER"] != "" {
		Postgres.User = myEnv["POSTGRES_USER"]
	}
	if myEnv["POSTGRES_PASSWORD"] != "" {
		Postgres.Password = myEnv["POSTGRES_PASSWORD"]
	}
	if myEnv["POSTGRES_DBNAME"] != "" {
		Postgres.DBName = myEnv["POSTGRES_DBNAME"]
	}
}
