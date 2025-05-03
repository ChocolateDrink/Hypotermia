package utils

import (
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
)

const (
	INPUT_KEYBOARD    = 1
	KEYEVENTF_KEYUP   = 0x0002
	KEYEVENTF_UNICODE = 0x0004
)

var sendInput *syscall.LazyProc = misc.User32.NewProc("SendInput")

type KeyInput struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type MouseInput struct {
	Dx        int32
	Dy        int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

type HardwareInput struct {
	UMsg    uint32
	WParamL uint16
	WParamH uint16
}

type Input struct {
	Type    uint32
	Ki      KeyInput
	Padding [8]byte
}

func SimulateInput(char rune) error {
	var inputs [2]Input

	inputs[0].Type = INPUT_KEYBOARD
	inputs[0].Ki = KeyInput{
		Scan:      uint16(char),
		Flags:     KEYEVENTF_UNICODE,
		Time:      0,
		ExtraInfo: 0,
	}

	inputs[1].Type = INPUT_KEYBOARD
	inputs[1].Ki = KeyInput{
		Scan:      uint16(char),
		Flags:     KEYEVENTF_UNICODE | KEYEVENTF_KEYUP,
		Time:      0,
		ExtraInfo: 0,
	}

	ret, _, err := sendInput.Call(
		uintptr(2),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret != 2 {
		return err
	}

	return nil
}
