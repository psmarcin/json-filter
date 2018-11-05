package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/itunes"
	"github.com/psmarcin/youtubeGoesPodcast/pkg/youtube"
)

func feedHandler(w http.ResponseWriter, r *http.Request) {
	youtubeURL := r.FormValue("youtubeUrl")
	log.SetPrefix("[FEED] ")
	defer log.SetPrefix("")

	log.Printf("Get %s %s", r.URL.RequestURI(), r.UserAgent())
	if youtubeURL == "" {
		err := errors.New("You need to provide channel id as query param 'channelId'")
		errorResponse(err, w)
		return
	}
	youtubeFeed, err := youtube.Create(youtubeURL)
	checkError(err, w)
	iTunesFeed := itunes.Create(youtubeFeed)
	xmlResponse(iTunesFeed.ToXML(), w)
}
