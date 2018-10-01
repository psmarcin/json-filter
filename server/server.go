package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/newrelic/go-agent"

	"github.com/psmarcin/youtubeGoesPodcast/itunes"
	agent "github.com/psmarcin/youtubeGoesPodcast/newrelic"
	"github.com/psmarcin/youtubeGoesPodcast/youtube"
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

func errorResponse(e error, w http.ResponseWriter) {
	err := Error{
		IsError:      true,
		Timestamp:    time.Now(),
		ErrorMessage: string(e.Error()),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resJSON, _ := json.Marshal(err)
	fmt.Fprint(w, string(resJSON))
}

func jsonResponse(b []byte, w http.ResponseWriter) {

	log.SetPrefix("[JSON] ")
	w.Header().Set("Content-Type", "application/json")
	s := string(b)
	log.Printf("Response %v", s)
	fmt.Fprintf(w, s)
}

func xmlResponse(b []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")
	s := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
` + string(b) + `
</rss>`
	fmt.Fprintf(w, s)
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.SetPrefix("[ROOT] ")
	resJSON, err := json.Marshal(rootStatus)
	if err != nil {
		errorResponse(err, w)
		return
	}
	jsonResponse(resJSON, w)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	channelID := r.FormValue("channelId")
	log.SetPrefix("[FEED] ")
	log.Printf("Request [%s] %s %s %s", r.Method, r.URL.RequestURI(), r.RemoteAddr, r.UserAgent())
	if channelID == "" {
		err := errors.New("You need to provide channel id as query param 'channelId'")
		errorResponse(err, w)
		return
	}
	youtubeFeed, err := youtube.Create(channelID)
	if err != nil {
		log.Print("Error ", err)
		errorResponse(err, w)
		return
	}
	iTunesFeed := itunes.Create(youtubeFeed)
	xmlResponse(iTunesFeed.ToXML(), w)
}

// Start creates server with fixed routes
func Start() {
	port := ":" + os.Getenv("PORT")
	a := agent.Init()
	log.SetPrefix("[SERVER] ")

	http.HandleFunc(newrelic.WrapHandleFunc(a, "/", handler))
	http.HandleFunc(newrelic.WrapHandleFunc(a, "/feed", feedHandler))

	log.Printf("Starting server at %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
