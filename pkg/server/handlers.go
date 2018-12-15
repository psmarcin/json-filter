package server

import (
	"encoding/json"
	"errors"
	"net/http"

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

	logger.Logger.Printf("%s %s %s %s %s %s", r.Method, r.URL.RequestURI(), r.UserAgent(), source, sourceType, contentType)
	if sourceType == "" || source == "" {
		err := errors.New("You need to provide source type and value")
		errorResponse(err, w, r)
		return
	}

	youtubeFeed, err := youtube.New(source, sourceType, querySearch)
	checkError(err, w, r)
	iTunesFeed := itunes.New(youtubeFeed)
	// if strings.Contains(contentType, "text/html") {
	// 	logger.Logger.Print("Response HTML")
	// 	htmlFeedResponse(iTunesFeed, w)
	// 	return
	// }
	xmlResponse(iTunesFeed.ToXML(), w)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	resJSON, err := json.Marshal(rootStatus)
	checkError(err, w, r)
	jsonResponse(resJSON, w)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger.Logger.Print("Get ", vars["videoId"], " ", r.Header)

	videoURL, err := getVideoURL(vars["videoId"])
	checkError(err, w, r)

	if videoURL.String() == "" {
		e := errors.New("Can't find proper source")
		checkError(e, w, r)
	}

	err = streamVideo(videoURL.String(), r.Header, w, r)
	checkError(err, w, r)
}
