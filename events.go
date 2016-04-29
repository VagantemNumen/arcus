package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ttacon/chalk"
)

type cmd struct {
	Cmd       string
	args      []string
	channelID string
}

var cmdChan = make(chan *cmd)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, "Arcus")
	state = s.State
	var cmds []Command
	cmds = append(cmds, ping)
	cmds = append(cmds, uptime)
	cmds = append(cmds, stats)
	go ProcessCommands(cmds)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := m.Message
	guild := getGuid(s, msg.ChannelID)
	channel := getChannel(s, msg.ChannelID)
	LogMessage(msg, guild, channel)

	prefix := config.App.Prefix

	if msg.Author.ID == arcus.ID || !strings.HasPrefix(msg.Content, prefix) || msg.Author.Bot {
		return
	}

	args := strings.Split(strings.TrimPrefix(strings.TrimSpace(msg.Content), prefix), " ")
	cmdText := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = nil
	}
	go func() {
		cmdChan <- &cmd{
			cmdText,
			args,
			msg.ChannelID,
		}
	}()
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	msg := m.Message
	guild := getGuid(s, msg.ChannelID)
	channel := getChannel(s, msg.ChannelID)
	fmt.Println(chalk.Magenta.Color(fmt.Sprintf("Message Updated at %v", msg.EditedTimestamp)))
	LogMessage(msg, guild, channel)
}

func getGuid(s *discordgo.Session, channelID string) *discordgo.Guild {
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
