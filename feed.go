package main

import (
	"fmt"
	"os"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
)

type FeedHandlers struct{}

func pollFeed(uri string, timeout int) {
	handlers := &FeedHandlers{}
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
	fmt.Printf("%d new rad item(s) in %s\n", len(newitems), feed.Url)
	/*for _, item := range newitems {
		fmt.Println("Title:", item.Title)
		for _, link := range item.Links {
			fmt.Println("Link:", link)
		}
		fmt.Println("Description:", item.Description)
		fmt.Println("Comments:", len(item.Comments))
		fmt.Println("Author:", item.Author)
		fmt.Println("PubDate:", item.PubDate)
	}*/
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
