package youtube

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/logger"
)

var youtubeChannelHttp = "http://www.youtube.com/channel/"
var youtubeChannelHttps = "https://www.youtube.com/channel/"
var channel = "https://www.googleapis.com/youtube/v3/channels"
var videos = "https://www.googleapis.com/youtube/v3/search"
var videoDetailsUrl = "https://www.googleapis.com/youtube/v3/videos"
var channelURL *url.URL
var videosURL *url.URL

var token = os.Getenv("PS_GOOGLE_API")

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

type Channel struct {
	Kind    string `json:"kind"`
	Etag    string `json:"etag"`
	ID      string `json:"id"`
	Snippet struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CustomURL   string `json:"customUrl"`
		PublishedAt string `json:"publishedAt"`
		Thumbnails  struct {
			Default struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		Localized struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"localized"`
		Country string `json:"country"`
	} `json:"snippet"`
}

type videoResponse struct {
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []Video `json:"items"`
}

type Video struct {
	Kind string `json:"kind"`
	Etag string `json:"etag"`
	ID   struct {
		Kind    string `json:"kind"`
		VideoID string `json:"videoId"`
	} `json:"id"`
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		ChannelID   string `json:"channelId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Thumbnails  struct {
			Default struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"default"`
			Medium struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"medium"`
			High struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"high"`
		} `json:"thumbnails"`
		ChannelTitle         string `json:"channelTitle"`
		LiveBroadcastContent string `json:"liveBroadcastContent"`
	} `json:"snippet"`
	Length int `json:"length"`
}

func (d *videoDetailsResponse) findById(id string) (videoDetails, error) {
	for _, v := range d.Items {
		if v.ID == id {
			return v, nil
			break
		}
	}
	return videoDetails{}, errors.New("Details not found")
}

func (yt *YouTube) getChannelURL() string {
	query := url.Values{}
	query.Add("key", os.Getenv("PS_GOOGLE_API"))
	query.Add("part", "snippet")
	query.Add("id", yt.ID)
	return channelURL.String() + "?" + query.Encode()
}

// GetChannel makes request to Google API and retreives snippet with basic information about channel
func (yt *YouTube) GetChannel() error {
	URL := yt.getChannelURL()
	response, err := http.Get(URL)
	if err != nil {
		logger.Logger.Fatal("Request ", err)
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Logger.Fatal("Read body ", err)
	}
	defer response.Body.Close()

	channel := channelResponse{}
	json.Unmarshal(content, &channel)
	if len(channel.Items) == 0 {
		return errors.New("Can't find channel")
	}
	yt.Channel = channel.Items[0]

	yt.Username = yt.Channel.Snippet.CustomURL

	// Unify publishedAt time
	publishedAt, err := time.Parse(time.RFC3339, yt.Channel.Snippet.PublishedAt)
	if err != nil {
		logger.Logger.Fatal("Parse publishedAt ", err)
	}
	yt.Channel.Snippet.PublishedAt = publishedAt.Format(time.RFC822)
	return nil
}

func (yt *YouTube) setChannelID(source, sourceType string) error {
	switch sourceType {
	case "channel":
		yt.ID = source
	case "channelUrl":
		parsedUrl, err := url.Parse(source)
		if err != nil {
			return err
		}
		split := strings.Split(parsedUrl.Path, "/")
		if len(split) < 3 {
			return errors.New("Wrong URL")
		}
		channelID := split[2] // channelId from parsedUrl
		yt.ID = channelID
	default:
		return errors.New("Source type not supported")
	}
	return nil
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
