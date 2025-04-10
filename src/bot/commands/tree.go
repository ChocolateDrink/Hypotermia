package commands

import (
	"fmt"
	"strconv"

	"Hypothermia/src/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	treeUsage string = "[path] [depth?]"

	treeArgsError string = "🟥 Expected 1 or more arguments."
	treeGenError  string = "🟥 Error in generating tree: "
	runError      string = "🟥 Error in running command: "
)

type TreeCommand struct{}

func (*TreeCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, treeArgsError+"\nUsage: "+treeUsage, m.Reference())
		return err
	}

	var depth int = 2
	var treeStr string

	if len(args) > 1 {
		num, err := strconv.Atoi(args[1])
		if err != nil {
			depth = 2
		} else {
			depth = num
		}
	}

	err := utils.GenerateTree(args[0], depth, 0, "", &treeStr)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(treeGenError+"%v", err), m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, treeStr, m.Reference())
	return err
}

func (*TreeCommand) Name() string {
	return "tree"
}

func (*TreeCommand) Info() string {
	return "shows a tree of files in a directory"
}
