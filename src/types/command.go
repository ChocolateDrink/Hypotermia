package types

import "github.com/bwmarrin/discordgo"

type Command interface {
	Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error
	Name() string
	Info() string
}
