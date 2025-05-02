package misc

var CmdShortcuts = map[string]string{
	"kill_tmgr":     "taskkill /IM Taskmgr.exe /F /T",
	"kill_dc":       "taskkill /IM discord.exe /F /T",
	"kill_settings": "taskkill /IM SystemSettings.exe /F /T",
	"shutdown":      "shutdown /s /f /t 0",
	"restart":       "shutdown /r /f /t 0",
	"logout":        "shutdown /l /f",
	"wifi":          "netsh wlan show profile",
	"sleep":         "rundll32.exe powrprof.dll,SetSuspendState 0,1,0",
	"ip":            "curl https://ipinfo.io/ip -s",
}

var PsShortcuts = map[string]string{
	"hwid":     "(Get-CimInstance Win32_ComputerSystemProduct).UUID",
	"kill_fg":  "Get-Process | Where-Object { $_.MainWindowHandle -ne 0 } | ForEach-Object { Stop-Process -Id $_.Id -Force }",
	"blackout": "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.Form]::new() | ForEach-Object {$_.WindowState='Maximized'; $_.FormBorderStyle='None'; $_.BackColor=[System.Drawing.Color]::Black; $_.TopMost=$true; $_.KeyPreview=$true; $_.Add_KeyDown({if($_.KeyCode -eq 'Escape'){$_.Close()}}); $_.ShowDialog()}",
	"clear_dt": "Get-ChildItem '$env:USERPROFILE\\Desktop' -Filter *.lnk | Remove-Item -Force",
}
