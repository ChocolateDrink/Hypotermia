package commands

import (
	"bytes"
	"image/jpeg"

	"github.com/bwmarrin/discordgo"
	"github.com/kbinani/screenshot"
)

const (
	ssCaptureError  string = "🟥 Failed to capture."
	ssEncodingError string = "🟥 Failed to encode screenshot."
)

type ScreenShotCommand struct{}

func (*ScreenShotCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.Capture(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
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
