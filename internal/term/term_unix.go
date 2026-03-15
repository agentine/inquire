//go:build !windows

package term

import (
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func tcgetattr(fd int, t *termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), ioctlGetTermios, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

func tcsetattr(fd int, t *termios) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), ioctlSetTermios, uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

func makeRaw(fd int) (*State, error) {
	var old termios
	if err := tcgetattr(fd, &old); err != nil {
		return nil, err
	}

	raw := old
	// Input: no break, no CR to NL, no parity, no strip, no flow control
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	// Output: disable post processing
	raw.Oflag &^= syscall.OPOST
	// Control: set 8-bit chars
	raw.Cflag |= syscall.CS8
	// Local: no echo, no canonical, no extended, no signal
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	// Read returns after 1 byte, no timeout
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	if err := tcsetattr(fd, &raw); err != nil {
		return nil, err
	}

	return &State{fd: fd, saved: old}, nil
}

func restore(s *State) error {
	old := s.saved.(termios)
	return tcsetattr(s.fd, &old)
}

func isTerminal(fd int) bool {
	var t termios
	return tcgetattr(fd, &t) == nil
}

func size(fd int) (int, int, error) {
	var ws winsize
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), ioctlGetWinsize, uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return 0, 0, errno
	}
	return int(ws.Col), int(ws.Row), nil
}
