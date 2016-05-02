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
	if err := session.ChannelTyping(channelID); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
	if _, err := session.ChannelMessageSend(channelID, res); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
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
	if err := session.ChannelTyping(channelID); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
	if _, err := session.ChannelMessageSend(channelID, res); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
}

func getMem() string {
	var m = new(runtime.MemStats)
	runtime.ReadMemStats(m)
	alloc := m.Alloc / 1e6
	tallloc := m.TotalAlloc / 1e6
	return fmt.Sprintf("%dMB / %dMB", alloc, tallloc)
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
	fmt.Println("From whoami.process()")
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
	ch, _ := session.State.Channel(channelID)
	g := ch.GuildID
	guild, _ := session.State.Guild(g)
	mem, _ := session.GuildMember(g, msg.Author.ID)
	var roles []string
	for _, role := range mem.Roles {
		for _, gr := range guild.Roles {
			if role == gr.ID {
				roles = append(roles, gr.Name)
			}
		}
	}
	joined, _ := time.Parse("2006-01-02T15:04:05.000000-07:00", mem.JoinedAt)
	res := "```rb\n"
	res += fmt.Sprintf("%-15s %s  '%s'\n", "Name", ":", strings.Replace(msg.Author.Username, "'", "’", -1))
	res += fmt.Sprintf("%-15s %s  '%s'\n", "Discriminator", ":", msg.Author.Discriminator)
	res += fmt.Sprintf("%-15s %s  '%s'\n", "ID", ":", msg.Author.ID)
	res += fmt.Sprintf("%-15s %s  '%s'\n", "Nickname", ":", mem.Nick)
	res += fmt.Sprintf("%-15s %s  '%t'\n", "Verified", ":", msg.Author.Verified)
	res += fmt.Sprintf("%-15s %s  '%v'\n", "Account Created", ":", ts.Format("January 02, 2006 15:04:05 MST"))
	res += fmt.Sprintf("%-15s %s  '%v'\n", "Joined At", ":", joined.UTC().Format("January 02, 2006 15:04:05 MST"))
	res += fmt.Sprintf("%-15s %s  '%s'\n", "Roles", ":", strings.Join(roles, ", "))
	res += "```"
	if err := session.ChannelTyping(channelID); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
	if _, err := session.ChannelMessageSend(channelID, res); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
	if _, err := session.ChannelFileSend(channelID, "avatar.jpg", avatar); err != nil {
		printError(fmt.Sprintf("%v", err))
	}
}

var whoami = Whoami{Name: "whoami"}

type Whois struct {
	Command
	Name string
}

func (c Whois) name() string {
	return c.Name
}

func (c Whois) process(channelID string, args []string, msg *discordgo.Message) {
	var res string
	var user *discordgo.User

	if len(msg.Mentions) > 0 {
		user = msg.Mentions[0]
	} else if len(args) > 0 && len(args[0]) >= 2 {
		for _, g := range session.State.Guilds {
			for _, mem := range g.Members {
				if strings.Contains(strings.ToUpper(mem.User.Username), strings.ToUpper(args[0])) {
					user = mem.User
				}
			}
		}
	}
	if user != nil {
		var ts time.Time
		if id, err := strconv.ParseInt(user.ID, 10, 64); err != nil {
			ts = time.Now()
		} else {
			ts = sf.Snowflake2utc(id)
		}
		url := discordgo.USER_AVATAR(user.ID, user.Avatar)
		response, _ := http.Get(url)
		defer response.Body.Close()
		avatar := response.Body
		ch, _ := session.State.Channel(channelID)
		g := ch.GuildID
		guild, _ := session.State.Guild(g)
		mem, _ := session.GuildMember(g, user.ID)
		var roles []string
		for _, role := range mem.Roles {
			for _, gr := range guild.Roles {
				if role == gr.ID {
					roles = append(roles, gr.Name)
				}
			}
		}
		joined, _ := time.Parse("2006-01-02T15:04:05.000000-07:00", mem.JoinedAt)
		res += "```rb\n"
		res += fmt.Sprintf("%-15s %s  '%s'\n", "Name", ":", strings.Replace(user.Username, "'", "’", -1))
		res += fmt.Sprintf("%-15s %s  '%s'\n", "Discriminator", ":", user.Discriminator)
		res += fmt.Sprintf("%-15s %s  '%s'\n", "ID", ":", user.ID)
		res += fmt.Sprintf("%-15s %s  '%s'\n", "Nickname", ":", mem.Nick)
		res += fmt.Sprintf("%-15s %s  '%v'\n", "Account Created", ":", ts.Format("January 02, 2006 15:04:05 MST"))
		res += fmt.Sprintf("%-15s %s  '%v'\n", "Joined At", ":", joined.UTC().Format("January 02, 2006 15:04:05 MST"))
		res += fmt.Sprintf("%-15s %s  '%s'\n", "Roles", ":", strings.Join(roles, ", "))
		res += "```"

		if err := session.ChannelTyping(channelID); err != nil {
			printError(fmt.Sprintf("%v", err))
		}
		if _, err := session.ChannelMessageSend(channelID, res); err != nil {
			printError(fmt.Sprintf("%v", err))
		}
		if _, err := session.ChannelFileSend(channelID, "avatar.jpg", avatar); err != nil {
			printError(fmt.Sprintf("%v", err))
		}
	}
}

var whois = Whois{Name: "whois"}
