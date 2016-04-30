package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	process(channelID string, args []string, msg *discordgo.Message)
	name() string
}

func ProcessCommands(cmds []Command) {
	for {
		cmd := <-cmdChan
		printInfo(fmt.Sprintf("Recieved command: %s, with args: %v", cmd.Cmd, cmd.args))
		go func() {
			for _, command := range cmds {
				if command.name() == cmd.Cmd {
					go command.process(cmd.channelID, cmd.args, cmd.msg)
					break
				}
			}
		}()
	}
}

func deleteMessage(duration int64, channelID string, messageID string) {
	time.Sleep(time.Duration(duration) * time.Second)
	session.ChannelMessageDelete(channelID, messageID)
}
