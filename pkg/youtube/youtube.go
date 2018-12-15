package youtube

import (
	"net/url"
	"os"
)

const youtubeChannelHTTP = "http://www.youtube.com/channel/"
const youtubeChannelHTTPS = "https://www.youtube.com/channel/"
const channel = "https://www.googleapis.com/youtube/v3/channels"
const videos = "https://www.googleapis.com/youtube/v3/search"
const videoDetailsURL = "https://www.googleapis.com/youtube/v3/videos"

var token = os.Getenv("PS_GOOGLE_API")

var channelURL *url.URL
var videosURL *url.URL

// YouTube is struct to store base information about channel and videos
type YouTube struct {
	ID       string
	Username string
	Channel  Channel
	Videos   []Video
	params   params
}

type params struct {
	search string
}

type channelResponse struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []Channel `json:"items"`
}

func (yt *YouTube) getVideoURL() string {
	query := url.Values{}
	query.Add("key", os.Getenv("PS_GOOGLE_API"))
	query.Add("part", "snippet")
	query.Add("channelId", yt.Channel.ID)
	query.Add("maxResults", "25")
	query.Add("order", "date")
	query.Add("q", yt.params.search)
	return videosURL.String() + "?" + query.Encode()
}

func init() {
	videosURL, _ = url.Parse(videos)
	channelURL, _ = url.Parse(channel)
}

// New makes new variable of type YouTube and gets all detaisls
func New(id, sourceType, search string) (YouTube, error) {
	yt := YouTube{
		params: params{
			search: search,
		},
	}
	yt.setChannelID(id, sourceType)
	err := yt.GetChannel()
	if err != nil {
		return yt, err
	}
	yt.GetVideos()

	return yt, nil
}
