package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ttacon/chalk"
)

type cmd struct {
	Cmd       string
	args      []string
	channelID string
	msg       *discordgo.Message
}

var cmdChan = make(chan *cmd)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, "Polling Feeds!")
	state = s.State
	var cmds []Command
	cmds = append(cmds, ping, uptime, stats, whoami, whois, guildinfo)
	// cmds = append(cmds, )
	// cmds = append(cmds, )
	// cmds = append(cmds, )
	cmds = append(cmds)
	go ProcessCommands(cmds)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := m.Message
	guild := getGuild(s, msg.ChannelID)
	channel := getChannel(s, msg.ChannelID)
	LogMessage(msg, guild, channel)

	prefix := config.App.Prefix

	if msg.Author.ID == arcus.ID || !strings.HasPrefix(msg.Content, prefix) || msg.Author.Bot {
		return
	}

	content := messageWithReplacedMentions(msg)

	args := strings.Split(strings.TrimPrefix(strings.TrimSpace(content), prefix), " ")
	cmdText := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = nil
	}
	deleteSpace(&args, "")
	deleteSpace(&args, " ")
	go func() {
		cmdChan <- &cmd{
			cmdText,
			args,
			msg.ChannelID,
			msg,
		}
	}()
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	msg := m.Message
	if msg.Author != nil {
		guild := getGuild(s, msg.ChannelID)
		channel := getChannel(s, msg.ChannelID)
		fmt.Println(chalk.Magenta.Color(fmt.Sprintf("Message Updated at %v", msg.EditedTimestamp)))
		LogMessage(msg, guild, channel)
	}
}

func getGuild(s *discordgo.Session, channelID string) *discordgo.Guild {
	channel := getChannel(s, channelID)

	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		printError(fmt.Sprintf("Error getting channel: %v", err))
		done <- true
	}

	return guild
}

func getChannel(s *discordgo.Session, channelID string) *discordgo.Channel {
	channel, err := s.State.Channel(channelID)
	if err != nil {
		printError(fmt.Sprintf("Error getting channel: %v", err))
		done <- true
	}

	return channel
}

//LogMessage logs messages to the console.
func LogMessage(msg *discordgo.Message, guild *discordgo.Guild, channel *discordgo.Channel) {
	guildStyle := chalk.Blue.NewStyle().WithTextStyle(chalk.Bold)
	channelStyle := chalk.Cyan.NewStyle()
	nameStyle := chalk.Green.NewStyle().WithTextStyle(chalk.Bold)
	timeStyle := chalk.Magenta.NewStyle().WithTextStyle(chalk.Italic)

	fmt.Println(
		guildStyle.Style(fmt.Sprintf("[%s]", guild.Name)),
		channelStyle.Style(fmt.Sprintf("[%s]", channel.Name)),
		timeStyle.Style(fmt.Sprintf("at %v\n", msg.Timestamp)),
		nameStyle.Style(fmt.Sprintf("[%s]", msg.Author.Username)),
		fmt.Sprintf("> %s", msg.ContentWithMentionsReplaced()),
	)
}

//LogMessage logs messages to the console.
func LogDelMessage(msg *discordgo.Message, guild *discordgo.Guild, channel *discordgo.Channel) {
	fmt.Println(chalk.Red.Color("Message Deleted:"))
	LogMessage(msg, guild, channel)
}

func messageWithReplacedMentions(m *discordgo.Message) string {
	if m.Mentions == nil {
		return m.Content
	}
	content := m.Content
	for _, user := range m.Mentions {
		content = strings.Replace(content, fmt.Sprintf("<@%s>", user.ID), "", -1)
		content = strings.Replace(content, fmt.Sprintf("<@!%s>", user.ID), "", -1)
	}

	roleIDRegex := regexp.MustCompile("<@&[0-9]*>")
	content = roleIDRegex.ReplaceAllStringFunc(content, func(str string) string {
		roleID := str[3 : len(str)-1]
		c, err := session.State.Channel(m.ChannelID)
		if err != nil {
			return str
		}
		g, err := session.State.Guild(c.GuildID)
		if err != nil {
			return str
		}
		for _, r := range g.Roles {
			if r.ID == roleID {
				return ""
			}
		}
		return str
	})

	channelIDRegex := regexp.MustCompile("<#[0-9]*>")
	content = channelIDRegex.ReplaceAllStringFunc(content, func(str string) string {
		_, err := session.State.Channel(str[2 : len(str)-1])
		if err != nil {
			return str
		}

		return ""
	})
	return content
}

func deleteSpace(s *[]string, selector string) {
	var temp []string
	for _, str := range *s {
		if str != selector {
			temp = append(temp, str)
		}
	}
	*s = temp
}
