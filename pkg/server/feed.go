package server

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/psmarcin/youtubeGoesPodcast/pkg/itunes"
	"github.com/psmarcin/youtubeGoesPodcast/pkg/youtube"
)

func baseHandler(w http.ResponseWriter, r *http.Request, sourceType, source, contentType string) {
	log.SetPrefix("[FEED] ")
	defer log.SetPrefix("")

	log.Printf("Get %s %s %s %s", r.URL.RequestURI(), r.UserAgent(), source, sourceType)
	if sourceType == "" || source == "" {
		err := errors.New("You need to provide channel id as query param 'channelId'")
		errorResponse(err, w)
		return
	}

	youtubeFeed, err := youtube.Create(source, sourceType)
	checkError(err, w)
	iTunesFeed := itunes.Create(youtubeFeed)

	// http request
	if strings.Contains(contentType, "text/html") {
		htmlFeedResponse(iTunesFeed, w)
		return
	}
	// default
	xmlResponse(iTunesFeed.ToXML(), w)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	value := r.FormValue("youtubeUrl")
	parsed, err := url.Parse(value)
	checkError(err, w)
	splitedPath := strings.Split(parsed.EscapedPath(), "/")
	if len(splitedPath) < 3 {
		err := errors.New("YouTube channel url not correct")
		checkError(err, w)
	}
	channelID := splitedPath[2] // channelId
	http.Redirect(w, r, "/feed/channel/"+channelID, http.StatusPermanentRedirect)
}

func feedPathHandler(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	baseHandler(w, r, parameters["sourceType"], parameters["source"], r.Header.Get("accept"))
}
