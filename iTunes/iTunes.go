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
	Category    string   `xml:"category"`
	Image       struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	Item []Item `xml:"item"`
}

type Item struct {
	GUID        string `xml:"guid"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL    string `xml:"url,attr"`
		Type   string `xml:"type,attr"`
		Length int    `xml:"length,attr"`
	} `xml:"enclosure"`
	Subtitle string `xml:"itines:subtitle"`
	Image    struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	Duration    string `xml:"itunes:duration"`
	Order       int    `xml:"itunes:order"`
	ITitle      string `xml:"itunes:title"`
	Summary     string `xml:"itunes:summary"`
	Author      string `xml:"itunes:author"`
	EpisodeType string `xml:"itunes:episodeType"`
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
		Author:   f.Title,
		Title:    f.Title,
		Link:     f.Link.Href,
		Category: "TV",
		Owner: Owner{
			Email: "p@pp.pp",
		},
		Language:    "en-us",
		PubDate:     f.Published,
		Description: f.ChannelDetails.Snippet.Description,
	}
	feed.Image.Href = f.ChannelDetails.Snippet.Thumbnails.High.URL
	items := []Item{}
	for i, v := range f.Feeds {
		item := Item{
			Title:       v.Title,
			ITitle:      v.Title,
			Subtitle:    v.Title,
			PubDate:     v.Published,
			Author:      v.Author,
			Link:        v.Link.Href,
			EpisodeType: "full",
			Duration:    "11:32",
			Order:       i,
			GUID:        "pspod://" + f.ChannelID + "/" + v.YTID,
			Description: v.Description,
			Summary:     v.Description,
		}
		item.Image.Href = "https://i.ytimg.com/vi/" + v.YTID + "/maxresdefault.jpg"
		item.Enclosure.URL = "http://podsync.net/download/PNyUU6D62/" + v.YTID + ".mp4?exp=tmp"
		item.Enclosure.Type = "video/mp4"
		item.Enclosure.Length = 242200000
		items = append(items, item)
	}
	feed.Item = items
	return feed
}
