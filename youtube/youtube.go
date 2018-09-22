package youtube

import (
	"encoding/json"
	"encoding/xml"

	feed "github.com/psmarcin/youtubeGoesPodcast/feed"
	iTunes "github.com/psmarcin/youtubeGoesPodcast/iTunes"
)

type Channel struct {
	ID     string      `json:"youTubeId" xml:"id"`
	Feed   feed.Feed   `json:"feed" xml:"youTube"`
	ITunes iTunes.Feed `json:"iTunes" xml:"iTunes"`
}

type Video struct {
	ID   string `json:"videoId"`
	Name string `json:"name"`
}

func CreateChannel(j string) Channel {
	user := Channel{
		ID:   j,
		Feed: feed.Feed{},
	}

	return user
}

func (channel *Channel) GetVideos() {
	feed := feed.Get(channel.ID)
	channel.Feed = feed
}

func (channel *Channel) GenerateITunes() {
	channel.ITunes = iTunes.Feed{
		Author: channel.Feed.Author,
		Link:   channel.Feed.Link.Href,
	}
	for _, v := range channel.Feed.Feeds {
		item := iTunes.Item{
			Title:   v.Title,
			ITitle:  v.Title,
			PubDate: v.Published,
			Author:  v.Author,
			Link:    v.Link.Href,
		}
		item.GetMedia()
		channel.ITunes.Item = append(channel.ITunes.Item, item)
	}
}

func (channel *Channel) ToJSON() ([]byte, error) {
	s, err := json.Marshal(channel)
	return s, err
}

func (channel *Channel) ToXML() ([]byte, error) {
	xml, err := xml.MarshalIndent(channel.ITunes, "  ", "    ")
	return xml, err
}
