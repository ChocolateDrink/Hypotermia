package commands

import (
	"bytes"
	"image/jpeg"

	"github.com/bwmarrin/discordgo"
	"github.com/vova616/screenshot"
)

const (
	ssCaptureError  string = "🟥 Failed to capture."
	ssEncodingError string = "🟥 Failed to encode screenshot."
)

type ScreenShotCommand struct{}

func (*ScreenShotCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, ssCaptureError, m.Reference())
		return err
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, ssEncodingError, m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: &buf,
		}},
	})

	return err
}

func (*ScreenShotCommand) Name() string {
	return "ss"
}

func (*ScreenShotCommand) Info() string {
	return "takes a screenshot"
}
