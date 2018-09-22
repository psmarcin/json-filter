package iTunes

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Feed struct {
	Author      string   `xml:"itunes:author"`
	Description string   `xml:"description"`
	Item        []Item   `xml:"item"`
	Language    string   `xml:"language"`
	Link        string   `xml:"link"`
	Owner       Owner    `xml:"itunes:owner"`
	Subtitle    string   `xml:"itunes:subtitle"`
	Title       string   `xml:"title"`
	Copyright   string   `xml:"copyright"`
	XMLName     xml.Name `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	ITitle      string `xml:"itunes:title"`
	Author      string `xml:"itunes:author"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Enclosure   struct {
		URL  string `xml:"url,attr"`
		Type string `xml:"type,attr"`
	} `xml:"enclosure"`
}

type Owner struct {
	Email string `xml:"itunes:email"`
}

func (i *Item) GetMedia() {
	URL := "https://www.audiotube.org/test?url=" + i.Link
	fmt.Println("[GET]", URL)
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	parsed := string(resp)
	fmt.Println("parsed", parsed)
	i.Enclosure.URL, _ = url.QueryUnescape(parsed)
	i.Enclosure.Type = "audio/mpeg"

}

func Create() Feed {
	feed := Feed{}
	return feed
}
