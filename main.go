package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ttacon/chalk"
)

// Client struct to store info about bot.
type Client struct {
	Username string
	ID       string
	Avatar   string
	Version  string
	Author   string
}

const Version = "1.2.0"

var (
	startTime = time.Now()
	config    Configuration
	session   *discordgo.Session
	state     *discordgo.State
	done      = make(chan bool)
	sigchan   = make(chan os.Signal, 1)
	arcus     Client
)

func main() {
	var err error
	printInfo(fmt.Sprintf("Application started at: %v", startTime))

	config = GetConfig()

	printInfo("Creating Discord Session.")
	if session, err = discordgo.New("Bot " + config.Auth.Token); err != nil {
		printError(fmt.Sprintf("Error creating Discord session: %v", err))
		done <- true
	}

	if config.App.Debug {
		session.LogLevel = discordgo.LogDebug
	}

	session.ShouldReconnectOnError = true

	session.AddHandler(ready)

	session.AddHandler(messageCreate)

	session.AddHandler(messageUpdate)

  if err = session.Open(); err != nil {
    printError(fmt.Sprintf("Error opening discordgo session: %v", err))
		sigchan <- syscall.SIGINT
  }


	err = getClient()
	if err != nil {
		printError(fmt.Sprintf("Error getting bot: %v", err))
		sigchan <- syscall.SIGINT
	}

	printInfo("Starting to poll feeds.")
	for _, feed := range config.Feeds {
		go pollFeed(feed.FeedUrl, feed.Timeout, feed.ChannelId)
	}

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	printInfo(fmt.Sprintf("Alloc: %v", m.Alloc/1000000))
	printInfo(fmt.Sprintf("TotalAlloc: %v", m.TotalAlloc/1000000))
	printInfo(fmt.Sprintf("Sys: %v", m.Sys/1000000))

	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigchan
		printInfo(fmt.Sprintf("Recieved signal:%v.", sig))
		printInfo("Cleaning up...")
		if err = session.Close(); err != nil {
			printError(fmt.Sprintf("Error cose connection: %v", err))
		}
		done <- true
	}()
	printInfo("Press Ctrl-C to exit.")
	<-done
	fmt.Println("Exiting...")
	os.Exit(0)
}

func printInfo(str string) {
	infoStyle := chalk.Cyan.NewStyle()
	infoStyle = infoStyle.WithTextStyle(chalk.Bold)
	fmt.Println(infoStyle.Style("[INFO]:"), str)
}

func printError(str string) {
	errorStyle := chalk.Red.NewStyle()
	errorStyle = errorStyle.WithTextStyle(chalk.Bold)
	fmt.Println(errorStyle.Style("[Error]:"), str)
}

func getClient() error {
	u, err := session.User("@me")
	if err != nil {
		return err
	}

	arcus.Username = u.Username
	arcus.ID = u.ID
	arcus.Avatar = u.Avatar
	arcus.Author = config.App.Author
	arcus.Version = config.App.Version
	return nil
}
