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

func SearchChannel(query string) ([]Channel, error) {
	req, err := http.NewRequest(http.MethodGet, videos, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("part", "snippet")
	req.Header.Add("q", query)
	req.Header.Add("key", token)
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	channels := channelResponse{}
	err = json.Unmarshal(body, &channels)
	if err != nil {
		return nil, err
	}
	return channels.Items, nil
}
