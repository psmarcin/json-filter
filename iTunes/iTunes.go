package iTunes

import (
	"encoding/xml"

	"github.com/psmarcin/youtubeGoesPodcast/youtube"
)

var YOUTUBE_VIDEO = "https://www.youtube.com/watch?v="
var YOUTUBE_CHANNEL = "https://www.youtube.com/channel/"

type Feed struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Category      string `xml:"category"`
	Generator     string `xml:"generator"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	PubDate       string `xml:"pubDate"`
	Image         struct {
		URL   string `xml:"url"`
		Title string `xml:"title"`
		Link  string `xml:"link"`
	} `xml:"image"`
	Subtitle    string `xml:"itunes:subtitle"`
	ITunesImage struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	Author  string   `xml:"itunes:author"`
	XMLName xml.Name `xml:"channel"`
	Item    []Item   `xml:"item"`
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
	Subtitle string `xml:"itunes:subtitle"`
	Image    struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	Order   int    `xml:"itunes:order"`
	ITitle  string `xml:"itunes:title"`
	Summary string `xml:"itunes:summary"`
	Author  string `xml:"itunes:author"`
}

type Owner struct {
	Email string `xml:"itunes:email"`
}

func (f *Feed) ToXML() []byte {
	b, _ := xml.MarshalIndent(f, "  ", "  ")
	return b
}

func Create(yt youtube.YouTube) Feed {
	feed := Feed{
		Title:       yt.Channel.Snippet.Title,
		Link:        YOUTUBE_CHANNEL + yt.Channel.ID,
		Description: yt.Channel.Snippet.Description,
		Category:    "TV",
		Author:      yt.Channel.Snippet.Title,
		Subtitle:    yt.Channel.Snippet.Title,
		Generator:   "psPodcast",
		Language:    yt.Channel.Snippet.Country,
		PubDate:     yt.Channel.Snippet.PublishedAt,
	}
	feed.Image.URL = yt.Channel.Snippet.Thumbnails.High.URL
	feed.Image.Title = yt.Channel.Snippet.Title
	feed.Image.Link = yt.Channel.ID
	feed.ITunesImage.Href = yt.Channel.Snippet.Thumbnails.High.URL
	items := []Item{}

	for i, v := range yt.Videos {
		item := Item{
			GUID:        "pspod://" + yt.Channel.ID + "/" + v.ID.VideoID,
			Title:       v.Snippet.Title,
			Link:        YOUTUBE_VIDEO + v.ID.VideoID,
			Description: v.Snippet.Description,
			PubDate:     v.Snippet.PublishedAt,
			Subtitle:    v.Snippet.Title,
			Order:       i,
			ITitle:      v.Snippet.Title,
			Summary:     v.Snippet.Description,
			Author:      yt.Channel.Snippet.Title,
		}
		item.Image.Href = "https://i.ytimg.com/vi/" + v.ID.VideoID + "/maxresdefault.jpg"

		item.Enclosure.URL = "http://podsync.net/download/PNyUU6D62/" + v.ID.VideoID + ".mp4?exp=tmp"
		item.Enclosure.Type = "video/mp4"
		item.Enclosure.Length = 242200000

		items = append(items, item)
	}
	feed.Item = items
	return feed
}
