package server

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/otium/ytdl"
)

// SOURCES are itag streams that we support, order count
var SOURCES = []interface{}{"139", "140", "141", "256", "258", "325", "328", "171", "172", "249", "250", "251", "5"}
var HEADER_FIELDS = []string{
	"Content-Length",
	"X-Content-Type-Options",
	"Last-Modified",
	"Cache-Control",
	"Connection",
	"Content-Type",
	"Date",
	"Expires",
	"Accept-Ranges",
	"Content-Range",
	"Range",
	"User-Agent",
}

func getVideoURL(videoID string) (*url.URL, error) {

	var videoURL *url.URL
	vid, err := ytdl.GetVideoInfoFromID(videoID)
	if err != nil {
		return nil, err
	}
	formats := vid.Formats.Filter(ytdl.FormatItagKey, SOURCES)
	if len(formats) == 0 {
		return nil, errors.New("Can't find proper source")
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
	return videoURL, nil
}

func setHeaders(src, dest http.Header) {
	for _, key := range HEADER_FIELDS {
		val := src.Get(key)
		if len(val) == 0 {
			continue
		}
		dest.Set(key, val)
	}
}

func streamVideo(url string, header http.Header, w http.ResponseWriter, method string) error {
	log.SetPrefix("[STREAM] ")
	defer log.SetPrefix("")

	req, err := http.NewRequest(method, url, nil)
	checkError(err, w)

	// set request header
	setHeaders(header, req.Header)

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	checkError(err, w)

	// set respoinse header
	setHeaders(resp.Header, w.Header())
	val := header.Get("Range")
	if val != "" {
		w.WriteHeader(http.StatusPartialContent)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	io.Copy(w, resp.Body)
	return nil
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[VIDEO GET] ")
	defer log.SetPrefix("")

	vars := mux.Vars(r)
	log.Print("Get ", vars["videoId"], " ", r.UserAgent())

	videoURL, err := getVideoURL(vars["videoId"])
	checkError(err, w)

	if videoURL.String() == "" {
		e := errors.New("Can't find proper source")
		checkError(e, w)
	}

	err = streamVideo(videoURL.String(), r.Header, w, http.MethodGet)
	checkError(err, w)
}

func videoHeadHandler(w http.ResponseWriter, r *http.Request) {
	log.SetPrefix("[VIDEO HEAD] ")
	defer log.SetPrefix("")

	vars := mux.Vars(r)
	log.Print("Get ", vars["videoId"], " ", r.UserAgent())

	videoURL, err := getVideoURL(vars["videoId"])
	checkError(err, w)

	if videoURL.String() == "" {
		e := errors.New("Can't find proper source")
		checkError(e, w)
	}

	err = streamVideo(videoURL.String(), r.Header, w, http.MethodHead)
	checkError(err, w)
}
