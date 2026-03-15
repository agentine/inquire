//go:build darwin

package term

// termios matches struct termios on Darwin (64-bit fields, NCCS=20).
type termios struct {
	Iflag  uint64
	Oflag  uint64
	Cflag  uint64
	Lflag  uint64
	Cc     [20]byte
	Ispeed uint64
	Ospeed uint64
}

const (
	ioctlGetTermios = 0x40487413 // TIOCGETA
	ioctlSetTermios = 0x80487414 // TIOCSETA
	ioctlGetWinsize = 0x40087468 // TIOCGWINSZ
)
