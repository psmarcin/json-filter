package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Init load env variables from files
func Init() {
	log.SetPrefix("[CONFIG] ")
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Can't load env variables ", err)
	}

	log.Print("PS_GOOGLE_API= ", os.Getenv("PS_GOOGLE_API"))
	log.Print("PS_NEWRELIC_KEY= ", os.Getenv("PS_NEWRELIC_KEY"))
	log.Print("PORT= ", os.Getenv("PORT"))
}
