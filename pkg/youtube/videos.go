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

type videoDetailsResponse struct {
	Items []videoDetails `json:"items"`
}

type videoDetails struct {
	ID             string         `json:"id"`
	ContentDetails contentDetails `json:"contentDetails"`
}

type contentDetails struct {
	Duration string `json:"duration"`
}

func (d *videoDetailsResponse) findByID(id string) (videoDetails, error) {
	for _, v := range d.Items {
		if v.ID == id {
			return v, nil
			break
		}
	}
	return videoDetails{}, errors.New("Details not found")
}

func (yt *YouTube) getVideoDetails() error {
	var videoIds []string
	// Unify publishedAt time
	for _, v := range yt.Videos {
		videoIds = append(videoIds, v.ID.VideoID)
	}
	// get video details
	detailsURL, _ := url.Parse(videoDetailsURL)
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
		details, _ := videosDetailsResp.findByID(v.ID.VideoID)
		trimmed := strings.Trim(details.ContentDetails.Duration, "PT")
		parsed := strings.ToLower(trimmed)
		if parsed == "" {
			continue
		}
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
