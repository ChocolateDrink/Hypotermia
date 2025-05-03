package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const (
	uploadFFError       string = "🟥 Failed to create form file."
	uploadCopyError     string = "🟥 Failed to copy file content."
	uploadCloseError    string = "🟥 Failed to close file."
	uploadHttpError     string = "🟥 Failed to create http request."
	uploadHttpSendError string = "🟥 Failed to send http request."
	uploadFailError     string = "🟥 Failed to upload file."
	uploadReadError     string = "🟥 Failed to read response body."

	uploadWriteFieldError string = "🟥 Failed to write field: %s"
)

func UploadFile(path string, file io.Reader) (string, string) {
	var body bytes.Buffer
	mpWriter := multipart.NewWriter(&body)

	writer, err := mpWriter.CreateFormFile("file", path)
	if err != nil {
		return "", uploadFFError
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return "", uploadCopyError
	}

	err = mpWriter.WriteField("expires", "1")
	if err != nil {
		return "", fmt.Sprintf(uploadWriteFieldError, "expires")
	}

	err = mpWriter.Close()
	if err != nil {
		return "", uploadCloseError
	}

	req, err := http.NewRequest("POST", "https://0x0.st", &body)
	if err != nil {
		return "", uploadHttpError
	}

	req.Header.Set("Content-Type", mpWriter.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", uploadHttpSendError
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "uploadReadError"
	}

	url := strings.TrimSpace(string(resBody))

	if !strings.HasPrefix(url, "http") {
		return "", uploadFailError
	}

	return url, ""
}
