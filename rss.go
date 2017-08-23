package babyrss

import (
	"encoding/xml"
	"fmt"
)

type Rss struct {
}

type Feed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (item Item) String() string {
	return fmt.Sprintf("title: %s, link: %s", item.Title, item.Link)
}

func (rss Rss) decode(data []byte) Feed {
	var feed Feed
	xml.Unmarshal(data, &feed)
	return feed
}
