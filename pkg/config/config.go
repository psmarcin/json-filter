package config

import (
	"github.com/psmarcin/youtubeGoesPodcast/pkg/logger"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Init load env variables from files
func Init() {
	// if load load developmemnt config
	if os.Getenv("NOW") == "" {
		logger.Logger.Print("Load development config")
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
