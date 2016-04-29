package main

import (
	"fmt"
	"strings"
	"time"
)

type Uptime struct {
	Command
	Name string
}

func (c Uptime) process(channelID string, args []string) {
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
