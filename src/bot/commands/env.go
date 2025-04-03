package commands

import (
	"Hypotermia/src/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	envGetError string = "🟥 Failed to get the user profile."
)

type EnvCommand struct{}

func (*EnvCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	info, err := utils.GetUserProfile()
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, envGetError, m.Reference())
		return err
	}

	var infoStr string
	infoStr += "Username: " + info.Username + "\n"
	infoStr += "Name: " + info.Name + "\n"
	infoStr += "Gid: " + info.Gid + "\n"
	infoStr += "Uid: " + info.Uid + "\n"
	infoStr += "HomeDir: " + info.HomeDir + "\n"

	_, err = s.ChannelMessageSendReply(m.ChannelID, infoStr, m.Reference())
	return err
}

func (*EnvCommand) Name() string {
	return "env"
}

func (*EnvCommand) Info() string {
	return "returns info about the users environment"
}
