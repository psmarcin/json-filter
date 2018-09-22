package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/psmarcin/youtubeGoesPodcast/feed"
	"github.com/psmarcin/youtubeGoesPodcast/iTunes"
)

var port = ":8080"

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
	w.Header().Set("Content-Type", "application/json")
	s := string(b)
	log.Printf("[Response] %v", s)
	fmt.Fprintf(w, s)
}

func xmlResponse(b []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")
	s := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
` + string(b)
	// log.Printf("[Response] %v", s)
	fmt.Fprintf(w, s)
}

func handler(w http.ResponseWriter, r *http.Request) {
	resJSON, err := json.Marshal(rootStatus)
	if err != nil {
		errorResponse(err, w)
		return
	}
	jsonResponse(resJSON, w)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	channelID, ok := r.Form["channelId"]
	if !ok {
		err := errors.New("You need to provide channel id as query param 'channelId'")
		errorResponse(err, w)
		return
	}

	youtubeFeed := feed.Create(channelID[0])
	iTunesFeed := iTunes.Create(youtubeFeed)
	xmlResponse(iTunesFeed.ToXML(), w)
}

func Start() {
	http.HandleFunc("/feed", feedHandler)
	http.HandleFunc("/", handler)

	log.Printf("Starting server at %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
