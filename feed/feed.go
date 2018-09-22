package feed

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type Feed struct {
	Author string  `xml:"author>name"`
	Link   Link    `xml:"link"`
	Feeds  []Entry `xml:"entry"`
}

type Entry struct {
	ID          string `xml:"id"`
	YTID        string `xml:"yt:videoId"`
	YTChannelID string `xml:"yt:channelId"`
	Title       string `xml:"title"`
	Link        Link   `xml:"link"`
	Author      string `xml:"author>name"`
	Published   string `xml:"published"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

func Get(channelId string) Feed {
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=" + channelId
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	feed, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	parsedFeed := parse(feed)
	return parsedFeed
}

func parse(s []byte) Feed {
	feed := Feed{}
	xml.Unmarshal([]byte(s), &feed)
	return feed
}
