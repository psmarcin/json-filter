package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Init load env variables from files
func Init() {
	log.SetPrefix("[CONFIG] ")
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Can't load env variables ", err)
	}
}
