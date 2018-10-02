package server

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent"

	"github.com/gorilla/mux"
	agent "github.com/psmarcin/youtubeGoesPodcast/newrelic"
)

type Status struct {
	Ok        bool      `json:"ok"`
	StartedAt time.Time `json:"startedAt"`
}

type Error struct {
	IsError      bool      `json:"isError"`
	Timestamp    time.Time `json:"timestamp"`
	ErrorMessage string    `json:"error"`
}

var rootStatus = Status{
	Ok:        true,
	StartedAt: time.Now(),
}

const ERR500 = "Internal Error"

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// Start creates server with fixed routes
func Start() {
	// Get env variable for port
	port := ":" + os.Getenv("PORT")
	// New Relic setup
	a := agent.Init()

	log.SetPrefix("[SERVER] ")
	defer log.SetPrefix("")

	router := mux.NewRouter()
	assets := http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))

	log.Print(assets)
	// routes
	router.HandleFunc(newrelic.WrapHandleFunc(a, "/", rootHandler)).Methods(http.MethodGet)
	router.HandleFunc(newrelic.WrapHandleFunc(a, "/stats", statsHandler)).Methods(http.MethodGet)
	router.HandleFunc(newrelic.WrapHandleFunc(a, "/feed", feedHandler)).Methods(http.MethodGet)
	// static assets
	http.Handle(newrelic.WrapHandle(a, "/assets/", assets))
	// mount router
	http.Handle("/", router)

	log.Printf("Starting server at %v", port)
	// listen server
	log.Fatal(http.ListenAndServe(port, nil))
}
