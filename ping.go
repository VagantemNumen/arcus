package main

type PingPong struct {
	Name string
}

func (c PingPong) process(channelID string, args []string) {
	session.ChannelMessageSend(channelID, "pong")
}

func (c PingPong) name() string {
	return c.Name
}

var ping = PingPong{Name: "ping"}
