package main

import "github.com/bwmarrin/discordgo"

type PingPong struct {
	Name string
}

func (c PingPong) process(channelID string, args []string, msg *discordgo.Message) {
	session.ChannelTyping(channelID)
	session.ChannelMessageSend(channelID, "pong")
}

func (c PingPong) name() string {
	return c.Name
}

var ping = PingPong{Name: "ping"}
