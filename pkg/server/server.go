package server

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	log.SetPrefix("[SERVER] ")
	defer log.SetPrefix("")

	router := mux.NewRouter()
	assets := prometheus.InstrumentHandler("/assets", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// routes
	router.HandleFunc("/", prometheus.InstrumentHandlerFunc("/", rootHandler)).Methods(http.MethodGet)
	router.HandleFunc("/stats", prometheus.InstrumentHandlerFunc("/stats", statsHandler)).Methods(http.MethodGet)
	router.Handle("/metrics", prometheus.InstrumentHandler("/metrics", promhttp.Handler())).Methods(http.MethodGet)
	// feed
	router.HandleFunc("/feed", prometheus.InstrumentHandlerFunc("/feed", feedHandler)).Methods(http.MethodGet)
	router.HandleFunc("/feed/{sourceType}/{source}", prometheus.InstrumentHandlerFunc("/feed", feedPathHandler)).Methods(http.MethodGet)
	// video
	router.HandleFunc("/video/{videoId}", prometheus.InstrumentHandlerFunc("/video/{videoId}", videoHandler)).Methods(http.MethodGet)
	router.HandleFunc("/video/{videoId}", prometheus.InstrumentHandlerFunc("/video/{videoId}", videoHeadHandler)).Methods(http.MethodHead)

	// static assets
	http.Handle("/assets/", assets)
	// mount router
	http.Handle("/", router)

	log.Printf("Starting server at %v", port)
	// listen server
	log.Fatal(http.ListenAndServe(port, nil))
}
