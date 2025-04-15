package utils

var CmdShortcuts = map[string]string{
	"kill_tmgr": "taskkill /F /IM Taskmgr.exe >nul 2>&1",
}

var PsShortcuts = map[string]string{
	"hwid": "(Get-CimInstance Win32_ComputerSystemProduct).UUID",
}
