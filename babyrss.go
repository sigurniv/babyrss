package babyrss

import (
	"net/http"
	"time"

	"errors"
	"io/ioutil"
	"log"
)

type RssStreamer struct {
	streamChan     chan Item
	url            string
	updateInterval time.Duration
	lastUpdateTime time.Time
}

var ErrFetch = errors.New("Error fetching rss")

func New(url string) *RssStreamer {
	return &RssStreamer{
		url:            url,
		streamChan:     make(chan Item),
		updateInterval: time.Second * 5,
		lastUpdateTime: time.Now(),
	}
}

func (streamer *RssStreamer) SetUpdateInterval(interval time.Duration) *RssStreamer {
	streamer.updateInterval = interval
	return streamer
}

func (streamer *RssStreamer) GetUpdatesChan() chan Item {
	go streamer.getUpdates()

	return streamer.streamChan
}

func (streamer *RssStreamer) getUpdates() {
	ticker := time.NewTicker(streamer.updateInterval)
	defer ticker.Stop()

	gettingUpdates := false

	rss := &Rss{}

	for {
		select {
		case <-ticker.C:
			if !gettingUpdates {
				log.Println("Getting rss updates")
				gettingUpdates = true

				body, err := streamer.fetch(streamer.url)
				if err != nil {
					log.Println(err)
					continue
				}

				feed := rss.decode(body)
				for _, item := range feed.Channel.Items {
					itemTime, err := ParseDate(item.PubDate)
					if err != nil {
						log.Println(err)
						continue
					}

					if !itemTime.After(streamer.lastUpdateTime) {
						continue
					} else {
						streamer.lastUpdateTime = itemTime
					}

					log.Printf("itemTime: %s, strTime: %s", itemTime, streamer.lastUpdateTime)
					streamer.streamChan <- item
				}

				gettingUpdates = false
			}
		}
	}
}

func (streamer *RssStreamer) fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, ErrFetch
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
