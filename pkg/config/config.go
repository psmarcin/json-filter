package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Init load env variables from files
func Init() {
	log.SetPrefix("[CONFIG] ")
	// if loacl load developmemnt config
	if os.Getenv("NOW") == "" {
		log.Print("Load development config")
		err := godotenv.Load(".env.dev")
		if err != nil {
			log.Print("Can't load .env.dev file ", err)
		}
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Print("Can't load .env file ", err)
		}
	}

	// default PORT value
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
}
