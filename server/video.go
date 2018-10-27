package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/otium/ytdl"
)

// SOURCES are itag streams that we support, order count
var SOURCES = []interface{}{"139", "140", "141", "256", "258", "325", "328", "171", "172", "249", "250", "251", "5"}

func streamVideo(url string, header http.Header, w http.ResponseWriter) error {
	log.SetPrefix("[STREAM] ")
	defer log.SetPrefix("")

	response, err := http.Get(url)
	checkError(err, w)

	log.Print("[STREAM] Headers ", header)
	response.Header.Set("Range", header.Get("Range"))

	w.Header().Set("Content-Length", response.Header["Content-Length"][0])
	w.Header().Set("X-Content-Type-Options", response.Header["X-Content-Type-Options"][0])
	w.Header().Set("Last-Modified", response.Header["Last-Modified"][0])
	w.Header().Set("Accept-Ranges", response.Header["Accept-Ranges"][0])
	w.Header().Set("Cache-Control", response.Header["Cache-Control"][0])
	w.Header().Set("Connection", response.Header["Connection"][0])
	w.Header().Set("Content-Type", response.Header["Content-Type"][0])
	w.Header().Set("Date", response.Header["Date"][0])
	w.Header().Set("Expires", response.Header["Expires"][0])
	io.Copy(w, response.Body)
	return nil
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[VIDEO] ")
	defer log.SetPrefix("")

	var videoURL *url.URL
	vars := mux.Vars(r)
	log.Print("Get ", vars["videoId"], " ", r.UserAgent())

	vid, err := ytdl.GetVideoInfoFromID(vars["videoId"])
	checkError(err, w)

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
		checkError(e, w)
	}

	err = streamVideo(videoURL.String(), r.Header, w)
	checkError(err, w)
}
