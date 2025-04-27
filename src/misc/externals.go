package misc

import "syscall"

var (
	User32 *syscall.LazyDLL = syscall.NewLazyDLL("user32.dll")
	NTdll  *syscall.LazyDLL = syscall.NewLazyDLL("ntdll.dll")

	AdjustPrivilege *syscall.LazyProc = NTdll.NewProc("RtlAdjustPrivilege")
)
