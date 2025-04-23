package utils

var CmdShortcuts = map[string]string{
	"kill_tmgr":     "taskkill /IM Taskmgr.exe /F /T",
	"kill_dc":       "taskkill /IM discord.exe /F /T",
	"kill_settings": "taskkill /IM SystemSettings.exe /F /T",
	"kill_all":      "for /f \"skip=3 tokens=1\" %i in ('tasklist /FI \"STATUS eq running\"') do taskkill /IM \"%i\" /F /T",
	"shutdown":      "shutdown /s /f /t 0",
	"restart":       "shutdown /r /f /t 0",
	"logout":        "shutdown /l /f",
	"sleep":         "rundll32.exe powrprof.dll,SetSuspendState 0,1,0",
	"ip":            "curl https://ipinfo.io/ip -s",
}

var PsShortcuts = map[string]string{
	"hwid":    "(Get-CimInstance Win32_ComputerSystemProduct).UUID",
	"kill_fg": "Get-Process | Where-Object { $_.MainWindowHandle -ne 0 } | ForEach-Object { Stop-Process -Id $_.Id -Force }",
}
