package youtube

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var youtubeChannelHttp = "http://www.youtube.com/channel/"
var youtubeChannelHttps = "https://www.youtube.com/channel/"
var channel = "https://www.googleapis.com/youtube/v3/channels"
var videos = "https://www.googleapis.com/youtube/v3/search"
var channelURL *url.URL
var videosURL *url.URL

var token = os.Getenv("PS_GOOGLE_API")

type YouTube struct {
	ID       string
	Username string
	Channel  Channel
	Videos   []Video
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
}

func (yt *YouTube) getVideoURL() string {
	query := url.Values{}
	query.Add("key", os.Getenv("PS_GOOGLE_API"))
	query.Add("part", "snippet")
	query.Add("channelId", yt.Channel.ID)
	query.Add("maxResults", "15")
	query.Add("order", "date")
	return videosURL.String() + "?" + query.Encode()
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
	log.SetPrefix("[YT CHANNEL] ")
	defer log.SetPrefix("")
	URL := yt.getChannelURL()
	log.Print("GET ", URL)
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal("Request ", err)
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Read body ", err)
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
		log.Fatal("Parse publishedAt ", err)
	}
	yt.Channel.Snippet.PublishedAt = publishedAt.Format(time.RFC1123Z)
	return nil
}

// GetVideos makes request to Google API and retreives last 15 videos snippets
func (yt *YouTube) GetVideos() {
	log.SetPrefix("[YT VIDEOS] ")
	defer log.SetPrefix("")

	URL := yt.getVideoURL()
	log.Print("GET ", URL)
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal("Request ", err)
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Read body ", err)
	}
	defer response.Body.Close()

	var videos videoResponse

	json.Unmarshal(content, &videos)

	// Unify publishedAt time
	for i, v := range videos.Items {
		pubAt, err := time.Parse(time.RFC3339, v.Snippet.PublishedAt)
		if err != nil {
			log.Fatal("Parse publishedAt ", err)
		}
		videos.Items[i].Snippet.PublishedAt = pubAt.Format(time.RFC1123Z)
	}

	yt.Videos = videos.Items
}

func init() {
	var err error
	videosURL, err = url.Parse(videos)
	if err != nil {
		log.Fatal("Cant' parse video url")
	}

	channelURL, err = url.Parse(channel)
	if err != nil {
		log.Fatal("Cant' parse channel url")
	}
}

// Create makes new variable of type YouTube and gets all detaisls
func Create(youtubeUrl string) (YouTube, error) {
	// TODO: Check if it's username or channelId
	log.Print("youtubeUrl", youtubeUrl)
	yt := YouTube{
		ID: youtubeUrl,
	}
	url, err := url.Parse(youtubeUrl)
	if err != nil {
		return yt, err
	}

	channelID := strings.Split(url.Path, "/")[2]
	yt.ID = channelID
	err = yt.GetChannel()
	if err != nil {
		return yt, err
	}
	yt.GetVideos()

	return yt, nil
}
