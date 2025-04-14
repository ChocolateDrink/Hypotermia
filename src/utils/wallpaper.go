package utils

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	wpGetPathError string = "🟥 Failed to get path."
	wpConvertError string = "🟥 Failed to convert argument."
	wpCantSetError string = "🟥 Failed to wallpaper."
)

const (
	SPI_SETDESKWALLPAPER int = 0x0014
	SPIF_UPDATEINIFILE   int = 0x01
	SPIF_SENDCHANGE      int = 0x02
)

var (
	user32       *syscall.LazyDLL  = syscall.NewLazyDLL("user32.dll")
	sysParamInfo *syscall.LazyProc = user32.NewProc("SystemParametersInfoW")
)

func SetWallpaper(file string) error {
	path, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf(wpGetPathError)
	}

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return fmt.Errorf(wpConvertError)
	}

	ret, _, _ := sysParamInfo.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		uintptr(0),
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	if ret == 0 {
		return fmt.Errorf(wpCantSetError)
	}

	return nil
}
