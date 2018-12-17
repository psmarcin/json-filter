package itunes

import (
	"encoding/xml"

	"github.com/psmarcin/youtubeGoesPodcast/pkg/youtube"
)

var youtubeVideoBaseURL = "https://www.youtube.com/watch?v="
var youtubeChannelBaseURL = "https://www.youtube.com/channel/"
var ygpProxyBaseURL = "http://youtube-goes-podcast-proxy.herokuapp.com/"

// Feed struct for JSON
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
	ITunesCategory struct {
		Text string `xml:"text,attr"`
	} `xml:"itunes:category"`
	Author  string   `xml:"itunes:author"`
	XMLName xml.Name `xml:"channel"`
	Item    []Item   `xml:"item"`
}

// Item struct for JSON
type Item struct {
	GUID        string    `xml:"guid"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     string    `xml:"pubDate"`
	Enclosure   Enclosure `xml:"enclosure"`
	Subtitle    string    `xml:"itunes:subtitle"`
	Image       struct {
		Href string `xml:"href,attr"`
	} `xml:"itunes:image"`
	Order    int    `xml:"itunes:order"`
	ITitle   string `xml:"itunes:title"`
	Summary  string `xml:"itunes:summary"`
	Author   string `xml:"itunes:author"`
	Duration int    `xml:"itunes:duration"`
}
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length int    `xml:"length,attr"`
}

// Owner struct
type Owner struct {
	Email string `xml:"itunes:email"`
}

// ToXML return XML
func (f *Feed) ToXML() []byte {
	b, _ := xml.MarshalIndent(f, "  ", "  ")
	return b
}

// New return new Feed
func New(yt youtube.YouTube) Feed {
	feed := Feed{
		Title:         yt.Channel.Snippet.Title,
		Link:          youtubeChannelBaseURL + yt.Channel.ID,
		Description:   yt.Channel.Snippet.Description,
		Category:      "TV &amp; Film",
		Author:        yt.Channel.Snippet.Title,
		Subtitle:      yt.Channel.Snippet.Title,
		Generator:     "psPodcast",
		Language:      yt.Channel.Snippet.Country,
		PubDate:       yt.Channel.Snippet.PublishedAt,
		LastBuildDate: yt.Channel.Snippet.PublishedAt,
	}
	feed.ITunesCategory.Text = "TV &amp; Film"
	feed.Image.URL = yt.Channel.Snippet.Thumbnails.High.URL
	feed.Image.Title = yt.Channel.Snippet.Title
	feed.Image.Link = yt.Channel.Snippet.Thumbnails.High.URL
	feed.ITunesImage.Href = yt.Channel.Snippet.Thumbnails.High.URL
	var items []Item

	for i, v := range yt.Videos {
		item := Item{
			GUID:        ygpProxyBaseURL + v.ID.VideoID,
			Title:       v.Snippet.Title,
			Link:        youtubeVideoBaseURL + v.ID.VideoID,
			Description: v.Snippet.Description,
			PubDate:     v.Snippet.PublishedAt,
			Subtitle:    v.Snippet.Title,
			Order:       i,
			ITitle:      v.Snippet.Title,
			Summary:     v.Snippet.Description,
			Author:      yt.Channel.Snippet.Title,
			Enclosure: Enclosure{
				URL:    ygpProxyBaseURL + v.ID.VideoID,
				Type:   "audio/mp3",
				Length: v.Length,
			},
		}
		item.Image.Href = "https://i.ytimg.com/vi/" + v.ID.VideoID + "/maxresdefault.jpg"
		item.Duration = v.Length
		items = append(items, item)
	}
	feed.Item = items
	return feed
}
