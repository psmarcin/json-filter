package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/psmarcin/youtubeGoesPodcast/itunes"
	"github.com/psmarcin/youtubeGoesPodcast/youtube"
)

func feedHandler(w http.ResponseWriter, r *http.Request) {
	youtubeURL := r.FormValue("youtubeUrl")
	log.SetPrefix("[FEED] ")
	defer log.SetPrefix("")

	log.Printf("Request [%s] %s %s %s", r.Method, r.URL.RequestURI(), r.RemoteAddr, r.UserAgent())
	if youtubeURL == "" {
		err := errors.New("You need to provide channel id as query param 'channelId'")
		errorResponse(err, w)
		return
	}
	youtubeFeed, err := youtube.Create(youtubeURL)
	if err != nil {
		log.Print("Error ", err)
		errorResponse(err, w)
		return
	}
	iTunesFeed := itunes.Create(youtubeFeed)
	xmlResponse(iTunesFeed.ToXML(), w)
}
