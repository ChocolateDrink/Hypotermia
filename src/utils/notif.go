package utils

const (
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_HELP              = 0x00004000
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_RETRYCANCEL       = 0x00000005
	MB_YESNO             = 0x00000004
	MB_YESNOCANCEL       = 0x00000003
)

const (
	ID_ABORT    = 3
	ID_CANCEL   = 2
	ID_CONTINUE = 11
	ID_IGNORE   = 5
	ID_NO       = 7
	ID_OK       = 1
	ID_RETRY    = 4
	ID_TRYAGAIN = 10
	ID_YES      = 6
)

func GetButtonFlag(button string) uint {
	switch button {
	case "abort-retry-ignore":
		return MB_ABORTRETRYIGNORE
	case "cancel-try_again-continue":
		return MB_CANCELTRYCONTINUE
	case "help":
		return MB_HELP
	case "ok":
		return MB_OK
	case "ok-cancel":
		return MB_OKCANCEL
	case "retry-cancel":
		return MB_RETRYCANCEL
	case "yes-no":
		return MB_YESNO
	case "yes-no-cancel":
		return MB_YESNOCANCEL
	default:
		return MB_OK
	}
}

func GetButtonClicked(button uintptr) string {
	switch int(button) {
	case ID_OK:
		return "OK"
	case ID_CANCEL:
		return "Cancel"
	case ID_ABORT:
		return "Abort"
	case ID_RETRY:
		return "Retry"
	case ID_IGNORE:
		return "Ignore"
	case ID_YES:
		return "Yes"
	case ID_NO:
		return "No"
	case ID_TRYAGAIN:
		return "Try Again"
	case ID_CONTINUE:
		return "Continue"
	default:
		return string(rune(button))
	}
}
