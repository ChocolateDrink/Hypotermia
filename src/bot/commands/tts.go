package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	ttsUsage string = "[...text]"

	ttsArgsError string = "🟥 Expected 1 argument."
	ttsRunError  string = "🟥 Failed to run TTS."

	ttsSuccess string = "🟩 Success in running TTS."
)

const (
	ttsAddSpeach string = "Add-Type -AssemblyName System.Speech"
	ttsSay       string = "(New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak(\"%s\")"
)

func (*TTSCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, ttsArgsError+"\nUsage: "+ttsUsage, m.Reference())
		return
	}

	cmdArg := fmt.Sprintf(ttsAddSpeach+"; "+ttsSay, strings.ReplaceAll(strings.Join(args, " "), "\"", "'"))
	cmd := exec.Command("powershell", "-Command", cmdArg)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	err := cmd.Run()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, ttsRunError, m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, ttsSuccess, m.Reference())
}

func (*TTSCommand) Name() string {
	return "tts"
}

func (*TTSCommand) Info() string {
	return "speaks the text with windows TTS"
}

type TTSCommand struct{}
