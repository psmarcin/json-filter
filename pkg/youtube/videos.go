package youtube

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/logger"
)

type videoDetailsResponse struct {
	Items []videoDetails `json:"items"`
}

type videoDetails struct {
	ID             string         `json:"id"`
	ContentDetails ContentDetails `json:"contentDetails"`
}

type ContentDetails struct {
	Duration string `json:"duration"`
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

// GetVideos makes request to Google API and retreives last 15 videos snippets
func (yt *YouTube) GetVideos() {
	URL := yt.getVideoURL()
	response, err := http.Get(URL)
	if err != nil {
		logger.Logger.Fatal("Request ", err)
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Logger.Fatal("Read body ", err)
	}
	defer response.Body.Close()

	var videos videoResponse
	json.Unmarshal(content, &videos)

	// Unify publishedAt time
	for i, v := range videos.Items {
		pubAt, err := time.Parse(time.RFC3339, v.Snippet.PublishedAt)
		if err != nil {
			logger.Logger.Fatal("Parse publishedAt ", err)
		}
		videos.Items[i].Snippet.PublishedAt = pubAt.Format(time.RFC822)
	}
	yt.Videos = videos.Items
	yt.getVideoDetails()
}

func (yt *YouTube) getVideoDetails() error {
	var videoIds []string
	// Unify publishedAt time
	for _, v := range yt.Videos {
		videoIds = append(videoIds, v.ID.VideoID)
	}
	// get video details
	detailsURL, _ := url.Parse(videoDetailsUrl)
	queryParams := detailsURL.Query()
	queryParams.Add("key", os.Getenv("PS_GOOGLE_API"))
	queryParams.Set("part", "contentDetails")
	queryParams.Set("fields", "items(contentDetails/duration,id)")
	queryParams.Set("id", strings.Join(videoIds, ","))

	detailsURL.RawQuery = queryParams.Encode()

	response, err := http.Get(detailsURL.String())
	if err != nil {
		logger.Logger.Fatal("Request ", err)
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Logger.Fatal("Read body ", err)
	}
	defer response.Body.Close()
	var videosDetailsResp videoDetailsResponse
	json.Unmarshal(content, &videosDetailsResp)

	for i, v := range yt.Videos {
		details, _ := videosDetailsResp.findById(v.ID.VideoID)
		trimmed := strings.Trim(details.ContentDetails.Duration, "PT")
		parsed := strings.ToLower(trimmed)
		l, err := time.ParseDuration(parsed)
		if err != nil {
			logger.Logger.Printf("Can't parse duration %s", err)
			return err
		}
		v.Length = int(l.Seconds())
		yt.Videos[i] = v
	}
	return nil
}
