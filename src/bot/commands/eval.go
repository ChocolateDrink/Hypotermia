package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	evalUsage string = "[cmd/powershell] [...args]"

	evalArgsError string = "🟥 Expected 2 or more arguments."
	evalUseError  string = "🟥 Invalid argument."
	evalRunError  string = "🟥 Error in running command: "

	evalNoOutput   string = "🟨 No output from command."
	evalGoodOutput string = "🟩 Success in running command."
	evalBadOutput  string = "🟥 Error in running command."
)

type EvalCommand struct{}

func (*EvalCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, evalArgsError+"\nUsage: "+evalUsage, m.Reference())
		return err
	}

	var cmd *exec.Cmd
	if args[0] == "powershell" || args[0] == "ps" {
		cmd = exec.Command("powershell", "-Command", strings.Join(args[1:], " "))
	} else if args[0] == "cmd" {
		cmd = exec.Command("cmd", "/C", strings.Join(args[1:], " "))
	} else {
		_, err := s.ChannelMessageSendReply(m.ChannelID, evalUseError+"\nUsage: "+evalUsage, m.Reference())
		return err
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		_, sendErr := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(evalRunError+"%v", err), m.Reference())
		if sendErr != nil {
			return sendErr
		}
	}

	response := ""

	output := out.String()
	if output != "" {
		response += evalGoodOutput + "\n```\n" + output + "\n```"
	}

	errorOutput := stderr.String()
	if errorOutput != "" {
		response += evalBadOutput + "\n```\n" + errorOutput + "\n```"
	}

	if response == "" {
		response = evalNoOutput
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
	return err
}

func (*EvalCommand) Name() string {
	return "eval"
}

func (*EvalCommand) Info() string {
	return "runs a command on cmd or powershell"
}
