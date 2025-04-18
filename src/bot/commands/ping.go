package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (*PingCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	msg := "pong"

	if len(args) > 0 {
		msg = "pong, args: " + strings.TrimSpace(strings.Join(args, " "))
	}

	s.ChannelMessageSend(m.ChannelID, msg)
}

func (*PingCommand) Name() string {
	return "ping"
}

func (*PingCommand) Info() string {
	return "test command"
}

type PingCommand struct{}
