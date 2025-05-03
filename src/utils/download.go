package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	dwCreateFileError string = "🟥 Failed to create file."
	dwHttpError       string = "🟥 Http error occurred."
	dwCopyError       string = "🟥 Failed to copy response body."
)

func DonwloadFile(url string, dest string) (string, string) {
	urlPath := strings.Split(url, "?")[0]
	fileName := filepath.Base(urlPath)

	var filePath string
	if dest == "" {
		filePath = filepath.Join(os.TempDir(), fileName)
	} else {
		info, err := os.Stat(dest)
		if err == nil && info.IsDir() {
			filePath = filepath.Join(dest, fileName)
		} else {
			filePath = dest
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", dwCreateFileError
	}

	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", dwHttpError
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", dwHttpError
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", fmt.Sprint(err)
	}

	return filePath, ""
}
