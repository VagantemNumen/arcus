package main

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
)

//FeedHandlers struct to hold the feed handlers as well as easy access to relevant fields.
type FeedHandlers struct {
	channelID  string
	lastUpdate time.Time
}

func pollFeed(uri string, timeout int, channelID string) {
	handlers := &FeedHandlers{
		channelID,
		time.Now().UTC(),
	}
	feed := rss.NewWithHandlers(timeout, true, handlers, handlers)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(time.Duration(30 * time.Second))
	}
}

func (fh *FeedHandlers) ProcessChannels(feed *rss.Feed, newchannels []*rss.Channel) {
	//fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (fh *FeedHandlers) ProcessItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	//fmt.Printf("%d new rad item(s) in %s\n", len(newitems), feed.Url)
	channelID := fh.channelID
	items := newitems
	regex := regexp.MustCompile(`<!--[\S\s]*?-->|<(?:"".*?""|'.*?'|[\S\s])*?>`)
	for i := len(items) - 1; i >= 0; i-- {
		pub, _ := items[i].ParsedPubDate()
		if pub.After(fh.lastUpdate) {
			var res string
			res += fmt.Sprintf("**%s**\n", items[i].Title)
			res += fmt.Sprintf("_**%s** - %s_\n", items[i].Author.Name, pub.UTC().Format("January 02, 2006 15:04:05 MST"))
			res += fmt.Sprintf("%s\n", html.UnescapeString(regex.ReplaceAllString(items[i].Description, "")))
			res += fmt.Sprintf("<%s>\n", items[i].Links[0].Href)
			if err := session.ChannelTyping(channelID); err != nil {
				printError(fmt.Sprintf("%v", err))
			}
			if _, err := session.ChannelMessageSend(channelID, res); err != nil {
				printError(fmt.Sprintf("%v", err))
			}
			fh.lastUpdate = pub
		}
	}
}
