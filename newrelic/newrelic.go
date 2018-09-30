package newrelic

import (
	"log"
	"os"

	newrelic "github.com/newrelic/go-agent"
)

// Init starts new newrelic instance and returns Application
func Init() newrelic.Application {
	log.SetPrefix("[NEW RELIC] ")
	name := "[DEV] podcast.psmarcin.me"
	if os.Getenv("NOW") != "" {
		name = "[PROD] podcast.psmarcin.me"
	}
	config := newrelic.NewConfig(name, os.Getenv("PS_NEWRELIC_KEY"))
	app, err := newrelic.NewApplication(config)
	if err != nil {
		log.Print("Error ", err)
		return nil
	}
	return app
}
