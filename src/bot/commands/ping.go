package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type PingCommand struct{}

func (*PingCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "pong, args: "+strings.TrimSpace(strings.Join(args, " ")))
	return err
}

func (*PingCommand) Name() string {
	return "ping"
}

func (*PingCommand) Info() string {
	return "test command"
}
