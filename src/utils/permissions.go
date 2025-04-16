package utils

import (
	"syscall"
	"unsafe"
)

const (
	SECURITY_BUILTIN_DOMAIN_RID uint32 = 0x20
	DOMAIN_ALIAS_RID_ADMINS     uint32 = 0x220
)

var (
	advapi32 *syscall.LazyDLL = syscall.NewLazyDLL("advapi32.dll")

	initSid    *syscall.LazyProc = advapi32.NewProc("AllocateAndInitializeSid")
	freeSid    *syscall.LazyProc = advapi32.NewProc("FreeSid")
	checkToken *syscall.LazyProc = advapi32.NewProc("CheckTokenMembership")
)

var identAuth [6]byte = [6]byte{0, 0, 0, 0, 0, 5}
var subAuth [2]uint32 = [2]uint32{SECURITY_BUILTIN_DOMAIN_RID, DOMAIN_ALIAS_RID_ADMINS}

func IsAdmin() (bool, error) {
	var sid uintptr
	var mem uint32

	ret, _, err := initSid.Call(
		uintptr(unsafe.Pointer(&identAuth[0])),
		uintptr(2),
		uintptr(subAuth[0]),
		uintptr(subAuth[1]),
		0, 0, 0, 0, 0, 0,
		uintptr(unsafe.Pointer(&sid)),
	)

	if ret == 0 {
		return false, err
	}

	defer freeSid.Call(sid)

	ret, _, err = checkToken.Call(
		0,
		sid,
		uintptr(unsafe.Pointer(&mem)),
	)

	if ret == 0 {
		return false, err
	}

	return mem != 0, nil
}

//func BypassUAC() error {
//
//}
