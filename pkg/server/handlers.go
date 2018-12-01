package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/psmarcin/youtubeGoesPodcast/pkg/itunes"
	"github.com/psmarcin/youtubeGoesPodcast/pkg/youtube"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.gohtml", nil)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)
	r.ParseForm()
	querySearch := r.FormValue("search")
	contentType := r.Header.Get("accept")
	source := string(parameters["source"])
	sourceType := string(parameters["sourceType"])

	logger.Logger.Printf("Get %s %s %s %s", r.URL.RequestURI(), r.UserAgent(), source, sourceType)
	if sourceType == "" || source == "" {
		err := errors.New("You need to provide source type and value")
		errorResponse(err, w)
		return
	}

	youtubeFeed, err := youtube.New(source, sourceType, querySearch)
	checkError(err, w)
	iTunesFeed := itunes.New(youtubeFeed)

	// http request
	if strings.Contains(contentType, "text/html") {
		htmlFeedResponse(iTunesFeed, w)
		return
	}
	// default
	xmlResponse(iTunesFeed.ToXML(), w)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	resJSON, err := json.Marshal(rootStatus)
	checkError(err, w)
	jsonResponse(resJSON, w)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger.Logger.Print("Get ", vars["videoId"], " ", r.Header)

	videoURL, err := getVideoURL(vars["videoId"])
	checkError(err, w)

	if videoURL.String() == "" {
		e := errors.New("Can't find proper source")
		checkError(e, w)
	}

	// logger.Logger.Printf("Request URL %s", videoURL.RequestURI())

	err = streamVideo(videoURL.String(), r.Header, w, http.MethodGet)
	checkError(err, w)
}

func videoHeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger.Logger.Print("Get ", vars["videoId"], " ", r.UserAgent())

	videoURL, err := getVideoURL(vars["videoId"])
	checkError(err, w)

	if videoURL.String() == "" {
		e := errors.New("Can't find proper source")
		checkError(e, w)
	}

	err = streamVideo(videoURL.String(), r.Header, w, http.MethodHead)
	checkError(err, w)
}
