package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	sf "github.com/VagantemNumen/arcus/discordsnowflake"
	"github.com/bwmarrin/discordgo"
)

// Uptime the struct for []uptime command.
type Uptime struct {
	Command
	Name string
}

func (c Uptime) process(channelID string, args []string, msg *discordgo.Message) {
	res := fmt.Sprintf("```xl\nUptime: %s\n```", getUptime())
	session.ChannelMessageSend(channelID, res)
}

func getUptime() string {
	uptime := time.Since(StartTime)
	t := int64(uptime.Nanoseconds() / 1000000)

	ms := t % 1000
	t /= 1000
	sec := t % 60
	t /= 60
	min := t % 60
	t /= 60
	hrs := t % 24
	days := t / 24

	var str string

	if days > 0 {
		str += fmt.Sprintf("%d %s", days, func() string {
			if days == 1 {
				return "day "
			} else {
				return "days "
			}
		}())
	}

	if hrs > 0 {
		str += fmt.Sprintf("%d %s", hrs, func() string {
			if hrs == 1 {
				return "hr "
			} else {
				return "hrs "
			}
		}())
	}

	if min > 0 {
		str += fmt.Sprintf("%d %s", min, func() string {
			if min == 1 {
				return "min "
			} else {
				return "mins "
			}
		}())
	}

	if sec > 0 {
		str += fmt.Sprintf("%d %s", sec, func() string {
			if sec == 1 {
				return "sec "
			} else {
				return "secs "
			}
		}())
	}

	if ms > 0 {
		str += fmt.Sprintf("%d ms", ms)
	}

	str = strings.TrimSpace(str) + "."
	return str
}

func (c Uptime) name() string {
	return c.Name
}

var uptime = Uptime{Name: "uptime"}

// Stats the struct for []stats command.
type Stats struct {
	Command
	Name string
}

func (c Stats) name() string {
	return c.Name
}

func (c Stats) process(channelID string, args []string, msg *discordgo.Message) {
	var channels []*discordgo.Channel
	guilds := session.State.Guilds
	for _, guild := range guilds {
		channels = append(channels, guild.Channels...)
	}
	res := "```rb\n"
	res += fmt.Sprintf("%-12s %s  '%s'\n", "Name", ":", arcus.Username)
	res += fmt.Sprintf("%-12s %s  '%s'\n", "ID", ":", arcus.ID)
	res += fmt.Sprintf("%-12s %s  '%s'\n", "Version", ":", arcus.Version)
	res += fmt.Sprintf("%-12s %s  '%d'\n", "Guilds", ":", len(guilds))
	res += fmt.Sprintf("%-12s %s  '%d'\n", "Channels", ":", len(channels))
	res += fmt.Sprintf("%-12s %s  '%s'\n", "Developer", ":", arcus.Author)
	res += fmt.Sprintf("%-12s %s  '%s'\n", "Go Version", ":", runtime.Version())
	res += fmt.Sprintf("%-12s %s  '%s'\n", "Memory Usage", ":", getMem())
	res += "```"
	session.ChannelMessageSend(channelID, res)
}

func getMem() string {
	var m = new(runtime.MemStats)
	runtime.ReadMemStats(m)
	alloc := m.Alloc / 1e6
	tallloc := m.TotalAlloc / 1e6
	sys := m.Sys / 1e6
	return fmt.Sprintf("%dMB / %dMB (%dMB)", alloc, tallloc, sys)
}

var stats = Stats{Name: "stats"}

// Whoami the struct for []whoami command.
type Whoami struct {
	Command
	Name string
}

func (c Whoami) name() string {
	return c.Name
}

func (c Whoami) process(channelID string, args []string, msg *discordgo.Message) {
	var ts time.Time
	if id, err := strconv.ParseInt(msg.Author.ID, 10, 64); err != nil {
		ts = time.Now()
	} else {
		ts = sf.Snowflake2utc(id)
	}

	url := discordgo.USER_AVATAR(msg.Author.ID, msg.Author.Avatar)
	response, _ := http.Get(url)

	defer response.Body.Close()
	avatar := response.Body

	res := "```rb\n"
	res += fmt.Sprintf("%-15s %s  '%s'\n", "Name", ":", msg.Author.Username)
	res += fmt.Sprintf("%-15s %s  '%s'\n", "Discriminator", ":", msg.Author.Discriminator)
	res += fmt.Sprintf("%-15s %s  '%s'\n", "ID", ":", msg.Author.ID)
	res += fmt.Sprintf("%-15s %s  '%t'\n", "Verfied", ":", msg.Author.Verified)
	res += fmt.Sprintf("%-15s %s  '%v'\n", "Account Created", ":", ts)
	res += "```"
	session.ChannelMessageSend(channelID, res)
	session.ChannelFileSend(channelID, "avatar.jpg", avatar)
}

var whoami = Whoami{Name: "whoami"}
