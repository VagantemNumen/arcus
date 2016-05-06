package main

import (
	"fmt"
	"html"
	"os"
	"regexp"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
)

var lastUpdate = time.Now()

type FeedHandlers struct {
	channelID string
}

func pollFeed(uri string, timeout int, channelID string) {
	handlers := &FeedHandlers{channelID}
	feed := rss.NewWithHandlers(timeout, true, handlers, handlers)

	for {
		if err := feed.Fetch(uri, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", uri, err)
			return
		}

		<-time.After(time.Duration(10 * time.Second))
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
		//if pub.After(lastUpdate) {
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
		lastUpdate = pub
		//}
	}
	/*
			// RSS and Shared fields
		    Title       string
		    Links       []*Link
		    Description string
		    Author      Author
		    Categories  []*Category
		    Comments    string
		    Enclosures  []*Enclosure
		    Guid        *string
		    PubDate     string
		    Source      *Source

		    // Atom specific fields
		    Id           string
		    Generator    *Generator
		    Contributors []string
		    Content      *Content
		    Updated      string
	*/
}
