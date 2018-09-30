package youTube

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var youTubeAPI = "https://www.googleapis.com/youtube/v3/channels?part=snippet&id="

type YouTube struct {
	ID string
}

type ChannelResp struct {
	Items []Channel `json:"items"`
}

type Channel struct {
	ID      string `json:"id"`
	Snippet ChannelSnippet
}

type ChannelSnippet struct {
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
}

var token = os.Getenv("PS_GOOGLE_API")

func (y *YouTube) GetChannel() Channel {
	log.SetPrefix("[GET CHANNEL] ")

	url := youTubeAPI + y.ID + "&key=" + token
	log.Printf("Get %s", y.ID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error ", err)
	}

	chResp := ChannelResp{}
	err = json.Unmarshal(body, &chResp)
	if err != nil {
		log.Fatal("Error ", err)
	}

	return chResp.Items[0]
}
