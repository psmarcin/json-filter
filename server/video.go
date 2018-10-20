package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/otium/ytdl"
)

// SOURCES are itag streams that we support, order count
var SOURCES = []interface{}{"139", "140", "141", "256", "258", "325", "328", "171", "172", "249", "250", "251", "5"}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[VIDEO] ")
	defer log.SetPrefix("")
	vars := mux.Vars(r)

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

	format := formats[0]
	u, err := vid.GetDownloadURL(format)
	if err != nil {
		errorResponse(err, w)
		return
	}
	w.Header().Set("location", u.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
