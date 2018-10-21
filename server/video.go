package server

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/otium/ytdl"
)

// SOURCES are itag streams that we support, order count
var SOURCES = []interface{}{"139", "140", "141", "256", "258", "325", "328", "171", "172", "249", "250", "251", "5"}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[VIDEO] ")
	defer log.SetPrefix("")

	var videoURL *url.URL
	vars := mux.Vars(r)
	log.Print("Get ", vars["videoId"], " ", r.UserAgent())

	vid, err := ytdl.GetVideoInfoFromID(vars["videoId"])
	if err != nil {
		errorResponse(err, w)
		return
	}

	formats := vid.Formats.Filter(ytdl.FormatItagKey, SOURCES)
	if len(formats) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse(errors.New("Can't find proper source"), w)
		return
	}

	// check urls if are valid (sometimes google return url that response with 403 status)
	for _, f := range formats {
		u, err := vid.GetDownloadURL(f)
		if err != nil {
			continue
		}
		_, err = http.Head(u.String())
		if err != nil {
			continue
		}
		videoURL = u
		break
	}
	if videoURL.String() == "" {
		e := errors.New("Can't find proper source")
		errorResponse(e, w)
		return
	}
	w.Header().Set("location", videoURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
