//go:build windows

package term

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode          = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode          = kernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

const (
	enableProcessedInput       = 0x0001
	enableLineInput            = 0x0002
	enableEchoInput            = 0x0004
	enableVirtualTermInput     = 0x0200
	enableProcessedOutput      = 0x0001
	enableVirtualTermProcessing = 0x0004
)

type consoleScreenBufferInfo struct {
	Size              [2]int16
	CursorPosition    [2]int16
	Attributes        uint16
	Window            [4]int16
	MaximumWindowSize [2]int16
}

func makeRaw(fd int) (*State, error) {
	h := syscall.Handle(fd)
	var mode uint32
	r, _, err := procGetConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return nil, err
	}

	raw := mode &^ (enableEchoInput | enableProcessedInput | enableLineInput)
	raw |= enableVirtualTermInput
	r, _, err = procSetConsoleMode.Call(uintptr(h), uintptr(raw))
	if r == 0 {
		return nil, err
	}

	return &State{fd: fd, saved: mode}, nil
}

func restore(s *State) error {
	h := syscall.Handle(s.fd)
	mode := s.saved.(uint32)
	r, _, err := procSetConsoleMode.Call(uintptr(h), uintptr(mode))
	if r == 0 {
		return err
	}
	return nil
}

func isTerminal(fd int) bool {
	h := syscall.Handle(fd)
	var mode uint32
	r, _, _ := procGetConsoleMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode)))
	return r != 0
}

func size(fd int) (int, int, error) {
	h := syscall.Handle(fd)
	var info consoleScreenBufferInfo
	r, _, err := procGetConsoleScreenBufInfo.Call(uintptr(h), uintptr(unsafe.Pointer(&info)))
	if r == 0 {
		return 0, 0, err
	}
	w := int(info.Window[2]-info.Window[0]) + 1
	ht := int(info.Window[3]-info.Window[1]) + 1
	return w, ht, nil
}
