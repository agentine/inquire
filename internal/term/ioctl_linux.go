//go:build linux

package term

// termios matches struct termios on Linux (32-bit fields, NCCS=32).
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   byte
	Cc     [32]byte
	Ispeed uint32
	Ospeed uint32
}

const (
	ioctlGetTermios = 0x5401 // TCGETS
	ioctlSetTermios = 0x5402 // TCSETS
	ioctlGetWinsize = 0x5413 // TIOCGWINSZ
)
