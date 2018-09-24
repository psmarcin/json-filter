package iTunes

import (
	"encoding/xml"

	"github.com/psmarcin/youtubeGoesPodcast/feed"
)

type Feed struct {
	Author      string   `xml:"itunes:author"`
	Description string   `xml:"description"`
	Language    string   `xml:"language"`
	Link        string   `xml:"link"`
	Owner       Owner    `xml:"itunes:owner"`
	Subtitle    string   `xml:"itunes:subtitle"`
	Title       string   `xml:"title"`
	Copyright   string   `xml:"copyright"`
	XMLName     xml.Name `xml:"channel"`
	PubDate     string   `xml:"pubDate"`
	Item        []Item   `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	ITitle      string `xml:"itunes:title"`
	Subtitle    string `xml:"itines:subtitle"`
	Description string `xml:"description"`
	Author      string `xml:"itunes:author"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL    string `xml:"url,attr"`
		Type   string `xml:"type,attr"`
		Length int    `xml:"length,attr"`
	} `xml:"enclosure"`
	Image struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
}

type Owner struct {
	Email string `xml:"itunes:email"`
}

func (f *Feed) ToXML() []byte {
	b, _ := xml.MarshalIndent(f, "  ", "    ")
	return b
}

func Create(f feed.Feed) Feed {
	feed := Feed{
		Author:  f.Title,
		Title:   f.Title,
		Link:    f.Link.Href,
		PubDate: f.Published,
	}
	items := []Item{}
	for _, v := range f.Feeds {
		item := Item{
			Title:       v.Title,
			ITitle:      v.Title,
			Subtitle:    v.Title,
			PubDate:     v.Published,
			Author:      v.Author,
			Link:        v.Link.Href,
			Description: v.Description,
		}
		item.Image.Href = "https://i.ytimg.com/vi/" + v.YTID + "/maxresdefault.jpg"
		item.Enclosure.URL = "http://podsync.net/download/PNyUU6D62/" + v.YTID + ".mp4"
		item.Enclosure.Type = "video/mp4"
		item.Enclosure.Length = 1000
		items = append(items, item)
	}
	feed.Item = items
	return feed
}
