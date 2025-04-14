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

func DonwloadFile(url string) (string, error) {
	urlPath := strings.Split(url, "?")[0]
	fileName := filepath.Base(urlPath)

	filePath := filepath.Join(os.TempDir(), fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf(dwCreateFileError)
	}

	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf(dwHttpError)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(dwHttpError)
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
