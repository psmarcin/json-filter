package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/itunes"
)

func checkError(e error, w http.ResponseWriter) {
	if e != nil {
		errorResponse(e, w)
		return
	}
}

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
	log.Printf("Response  %v", s)
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

func htmlFeedResponse(feed itunes.Feed, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	templates.ExecuteTemplate(w, "feed.gohtml", feed)
}
