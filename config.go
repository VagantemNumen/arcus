package main

import (
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

//Configuration is a struct used to store the configuration for the bot.
type Configuration struct {
	App struct {
		Name    string
		Version string
		Author  string
		Owner   string
		Prefix  string
		Debug   bool
	}
	Auth struct {
		Token string
	}
	Feeds []Feed
}

//Server definition
type Feed struct {
	FeedUrl   string
	ChannelId string
	Timeout   int
}

// Reload reloads the configuration in memory.
func (c *Configuration) Reload() {
	config := GetConfig()
	c = &config
}

// GetConfig Generates config from the provided "config.toml" file and returns a Configuration struct.
func GetConfig() Configuration {
	var config Configuration

	f, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}

	config.App.Version = Version

	return config
}
