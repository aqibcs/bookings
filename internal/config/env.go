package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Username     = "myuser"
	Password     = ""
	DatabaseName = "mydatabase"
	Host         = "localhost"
	Port         = "5432"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

	if x := os.Getenv("DB_USER"); x != "" {
		Username = x
	}

	if x := os.Getenv("DB_PASSWORD"); x != "" {
		Password = x
	}

	if x := os.Getenv("DB_NAME"); x != "" {
		DatabaseName = x
	}

	if x := os.Getenv("DB_HOST"); x != "" {
		Host = x
	}

	if x := os.Getenv("DB_PORT"); x != "" {
		Port = x
	}
}
