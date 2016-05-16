# Arcus

[![Build Status](https://travis-ci.org/VagantemNumen/arcus.svg?branch=master)](https://travis-ci.org/VagantemNumen/arcus) [![Build Status](https://drone.io/github.com/VagantemNumen/arcus/status.png)](https://drone.io/github.com/VagantemNumen/arcus/latest)

A Discord bot written in Go using discordgo Discord API wrapper.

## Function

The bot will poll for feeds provided in the config.toml file and post new feeds it sees into a channel specified for the feed.

### Additional Functions

```
- []whoami : Shows the user information of the person.
- []whois <part/full username>/<@mention> : Shows the user information of the matching person the bot finds.
- []guildinfo : Shows information about the guild.
- []uptime : Shows the uptime of the bot.
```

## Getting the release

Download the latest release from [Release](https://github.com/VagantemNumen/arcus/releases/latest).

Only available for amd64 Linux at the moment.

## Usage

Create a file named **config.toml** inside the current working directory. Optionally dowload config.toml.example from the repository and rename to config.toml.
Inside the file put the following content and replace the specified parts.

```toml
[app]
	name 	= "Arcus"
	version = "1.1.1"
	author 	= "AI"
	owner 	= "148793367126474752"
	prefix 	= "[]" #Change to the prefix you want to use for commands.
	debug 	= false  #Set to true to get additional output.

[auth]
	token 	= "YOUR_BOT_TOKEN"

#Add as many [[feeds]] blocks as needed in the provided format. Keep timeout at 0 for faster feeds.
[[feeds]]
	feed_url 	= "https://blog.discordapp.com/rss/"
	channel_id 	= "145645642868452776"
	timeout 	= 0
[[feeds]]
	feed_url 	= "https://blog.discordapp.com/rss/"
	channel_id 	= "145645642868452546"
	timeout 	= 0
```

Run the bot from your shell. If the bot is not in your path keep the executable inside the current working directory along with config.toml.
Rename bot to **arcus** for ease of use.

```sh
arcus
```

That is all. Enjoy your feeds.