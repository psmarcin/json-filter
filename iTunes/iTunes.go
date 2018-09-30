package iTunes

import (
	"encoding/xml"
	"log"
	"time"

	"github.com/psmarcin/youtubeGoesPodcast/feed"
)

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

func Create(f feed.Feed) Feed {
	pub, err := time.Parse(time.RFC3339, f.Published)
	if err != nil {
		log.Fatal("Time parse ", err)
	}
	feed := Feed{
		Title:       f.Title,
		Link:        f.Link.Href,
		Description: f.ChannelDetails.Snippet.Description,
		Category:    "TV",
		Author:      f.Title,
		Subtitle:    f.Title,
		Generator:   "psPodcast",
		Language:    "en-us",
		PubDate:     pub.Format(time.RFC1123Z),
	}
	feed.Image.URL = f.ChannelDetails.Snippet.Thumbnails.High.URL
	feed.Image.Title = f.Title
	feed.Image.Link = f.Link.Href
	feed.ITunesImage.Href = f.ChannelDetails.Snippet.Thumbnails.High.URL
	items := []Item{}

	for i, v := range f.Feeds {
		pub, err := time.Parse(time.RFC3339, v.Published)
		if err != nil {
			log.Fatal("Time parse ", err)
		}
		item := Item{
			GUID:        "pspod://" + f.ChannelID + "/" + v.YTID,
			Title:       v.Title,
			Link:        v.Link.Href,
			Description: v.Description,
			PubDate:     pub.Format(time.RFC1123Z),
			Subtitle:    v.Title,
			Order:       i,
			ITitle:      v.Title,
			Summary:     v.Description,
			Author:      v.Author,
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
