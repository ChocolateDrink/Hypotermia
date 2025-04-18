package utils

import "syscall"

var (
	User32 *syscall.LazyDLL = syscall.NewLazyDLL("user32.dll")
)
