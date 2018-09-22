package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	youtube "github.com/psmarcin/youtubeGoesPodcast/youtube"
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
<?xml-stylesheet type="text/xsl" media="screen" href="/~d/styles/rss2enclosuresfull.xsl"?><?xml-stylesheet type="text/css" media="screen" href="http://feeds.serialpodcast.org/~d/styles/itemcontent.css"?><rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:feedburner="http://rssnamespace.org/feedburner/ext/1.0" version="2.0" xml:base="https://serialpodcast.org">
` + string(b)
	// log.Printf("[Response] %v", s)
	fmt.Fprintf(w, s)
}

func handler(w http.ResponseWriter, r *http.Request) {
	resJSON, err := json.Marshal(rootStatus)
	if err != nil {
		errorResponse(err, w)
	}

	w.WriteHeader(http.StatusOK)
	jsonResponse(resJSON, w)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	channelID, ok := r.Form["channelId"]
	if ok != true {
		err := errors.New("Please provide channel id as 'channelId' query string")
		errorResponse(err, w)
		return
	}

	channel := youtube.CreateChannel(channelID[0])
	channel.GetVideos()
	channel.GenerateITunes()
	xml, err := channel.ToXML()
	if err != nil {
		errorResponse(err, w)
	}
	xmlResponse(xml, w)
}

func Start() {
	http.HandleFunc("/feed", feedHandler)
	http.HandleFunc("/", handler)

	log.Printf("Starting server at %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
