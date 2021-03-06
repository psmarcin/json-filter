package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/logger"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/itunes"
)

// status struct
type status struct {
	Ok        bool      `json:"ok"`
	StartedAt time.Time `json:"startedAt"`
}

type jsonError struct {
	IsError      bool      `json:"isError"`
	Timestamp    time.Time `json:"timestamp"`
	ErrorMessage string    `json:"error"`
}

var rootStatus = status{
	Ok:        true,
	StartedAt: time.Now(),
}

func checkError(e error, w http.ResponseWriter, r *http.Request) {
	if e != nil {
		errorResponse(e, w, r)
		return
	}
}

func errorResponse(e error, w http.ResponseWriter, r *http.Request) {
	err := jsonError{
		IsError:      true,
		Timestamp:    time.Now(),
		ErrorMessage: string(e.Error()),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	logger.Logger.Print("[ERROR] ", err)
	resJSON, _ := json.Marshal(err)
	fmt.Fprint(w, string(resJSON))
}

func jsonResponse(b []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	s := string(b)
	fmt.Fprintf(w, s)
}

func xmlResponse(b []byte, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/rss+xml; charset=UTF-8")
	s := `<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
` + string(b) + `
</rss>`
	fmt.Fprintf(w, s)
}

func htmlFeedResponse(feed itunes.Feed, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	templates.ExecuteTemplate(w, "feed.gohtml", feed)
}
