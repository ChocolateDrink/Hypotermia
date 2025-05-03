package utils

import (
	"path/filepath"
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
)

const (
	wpGetPathError string = "🟥 Failed to get path."
	wpCantSetError string = "🟥 Failed to wallpaper."
)

const (
	SPI_SETDESKWALLPAPER int = 0x0014
	SPIF_UPDATEINIFILE   int = 0x01
	SPIF_SENDCHANGE      int = 0x02
)

var sysParamInfo *syscall.LazyProc = misc.User32.NewProc("SystemParametersInfoW")

func SetWallpaper(file string) string {
	path, err := filepath.Abs(file)
	if err != nil {
		return wpGetPathError
	}

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return misc.ERROR_CONVERT
	}

	ret, _, _ := sysParamInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		uintptr(0),
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	if ret == 0 {
		return wpCantSetError
	}

	return ""
}
