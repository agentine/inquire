// Package keys provides key event parsing from raw terminal input.
package keys

import "io"

// Key represents a parsed key event.
type Key struct {
	Type KeyType
	Rune rune // populated for KeyRune
}

// KeyType identifies the type of key event.
type KeyType int

const (
	KeyRune      KeyType = iota // printable character
	KeyEnter                    // Enter/Return
	KeyBackspace                // Backspace
	KeyDelete                   // Delete
	KeyTab                      // Tab
	KeyEscape                   // Escape (standalone)
	KeyUp                       // Arrow up
	KeyDown                     // Arrow down
	KeyLeft                     // Arrow left
	KeyRight                    // Arrow right
	KeyHome                     // Home
	KeyEnd                      // End
	KeyPageUp                   // Page up
	KeyPageDown                 // Page down
	KeyCtrlA                    // Ctrl+A
	KeyCtrlB                    // Ctrl+B
	KeyCtrlC                    // Ctrl+C (interrupt)
	KeyCtrlD                    // Ctrl+D (EOF)
	KeyCtrlE                    // Ctrl+E
	KeyCtrlF                    // Ctrl+F
	KeyCtrlH                    // Ctrl+H (backspace alias)
	KeyCtrlK                    // Ctrl+K
	KeyCtrlL                    // Ctrl+L
	KeyCtrlN                    // Ctrl+N
	KeyCtrlP                    // Ctrl+P
	KeyCtrlU                    // Ctrl+U
	KeyCtrlW                    // Ctrl+W
	KeySpace                    // Space
)

// String returns a human-readable name for the key type.
func (k KeyType) String() string {
	switch k {
	case KeyRune:
		return "Rune"
	case KeyEnter:
		return "Enter"
	case KeyBackspace:
		return "Backspace"
	case KeyDelete:
		return "Delete"
	case KeyTab:
		return "Tab"
	case KeyEscape:
		return "Escape"
	case KeyUp:
		return "Up"
	case KeyDown:
		return "Down"
	case KeyLeft:
		return "Left"
	case KeyRight:
		return "Right"
	case KeyHome:
		return "Home"
	case KeyEnd:
		return "End"
	case KeyPageUp:
		return "PageUp"
	case KeyPageDown:
		return "PageDown"
	case KeyCtrlC:
		return "Ctrl+C"
	case KeyCtrlD:
		return "Ctrl+D"
	case KeySpace:
		return "Space"
	default:
		return "Unknown"
	}
}

// Reader reads raw bytes from a terminal and parses key events.
type Reader struct {
	in  io.Reader
	buf [256]byte
	pos int
	n   int
}

// NewReader creates a key reader from the given input.
func NewReader(in io.Reader) *Reader {
	return &Reader{in: in}
}

// ReadKey reads and returns the next key event.
func (r *Reader) ReadKey() (Key, error) {
	b, err := r.readByte()
	if err != nil {
		return Key{}, err
	}

	switch {
	case b == 0x1b: // Escape
		return r.parseEscape()
	case b == 0x0d || b == 0x0a: // CR or LF
		return Key{Type: KeyEnter}, nil
	case b == 0x7f: // DEL
		return Key{Type: KeyBackspace}, nil
	case b == 0x08: // Ctrl+H / Backspace
		return Key{Type: KeyCtrlH}, nil
	case b == 0x09: // Tab
		return Key{Type: KeyTab}, nil
	case b == 0x01:
		return Key{Type: KeyCtrlA}, nil
	case b == 0x02:
		return Key{Type: KeyCtrlB}, nil
	case b == 0x03:
		return Key{Type: KeyCtrlC}, nil
	case b == 0x04:
		return Key{Type: KeyCtrlD}, nil
	case b == 0x05:
		return Key{Type: KeyCtrlE}, nil
	case b == 0x06:
		return Key{Type: KeyCtrlF}, nil
	case b == 0x0b:
		return Key{Type: KeyCtrlK}, nil
	case b == 0x0c:
		return Key{Type: KeyCtrlL}, nil
	case b == 0x0e:
		return Key{Type: KeyCtrlN}, nil
	case b == 0x10:
		return Key{Type: KeyCtrlP}, nil
	case b == 0x15:
		return Key{Type: KeyCtrlU}, nil
	case b == 0x17:
		return Key{Type: KeyCtrlW}, nil
	case b == ' ':
		return Key{Type: KeySpace, Rune: ' '}, nil
	case b >= 0x20 && b < 0x7f: // printable ASCII
		return Key{Type: KeyRune, Rune: rune(b)}, nil
	case b >= 0xc0: // UTF-8 multi-byte
		return r.parseUTF8(b)
	default:
		return Key{Type: KeyRune, Rune: rune(b)}, nil
	}
}

func (r *Reader) parseEscape() (Key, error) {
	// Try to read the next byte with a quick check
	b, err := r.peekByte()
	if err != nil {
		// Standalone escape
		return Key{Type: KeyEscape}, nil
	}

	if b != '[' && b != 'O' {
		// Alt+key or standalone escape — consume and return escape
		return Key{Type: KeyEscape}, nil
	}

	// Consume the [ or O
	r.consumeByte()

	if b == 'O' {
		// SS3 sequences (e.g., OF = End, OH = Home)
		next, err := r.readByte()
		if err != nil {
			return Key{Type: KeyEscape}, nil
		}
		switch next {
		case 'F':
			return Key{Type: KeyEnd}, nil
		case 'H':
			return Key{Type: KeyHome}, nil
		default:
			return Key{Type: KeyEscape}, nil
		}
	}

	// CSI sequences: ESC [ ...
	next, err := r.readByte()
	if err != nil {
		return Key{Type: KeyEscape}, nil
	}

	switch next {
	case 'A':
		return Key{Type: KeyUp}, nil
	case 'B':
		return Key{Type: KeyDown}, nil
	case 'C':
		return Key{Type: KeyRight}, nil
	case 'D':
		return Key{Type: KeyLeft}, nil
	case 'H':
		return Key{Type: KeyHome}, nil
	case 'F':
		return Key{Type: KeyEnd}, nil
	case '1', '2', '3', '4', '5', '6', '7', '8':
		return r.parseCSIParam(next)
	default:
		return Key{Type: KeyEscape}, nil
	}
}

func (r *Reader) parseCSIParam(first byte) (Key, error) {
	// Read until we get a terminator (~ or a letter)
	param := []byte{first}
	for {
		b, err := r.readByte()
		if err != nil {
			return Key{Type: KeyEscape}, nil
		}
		if b == '~' {
			break
		}
		if b >= 'A' && b <= 'Z' {
			// CSI 1;5C style sequences
			return Key{Type: KeyEscape}, nil
		}
		param = append(param, b)
	}

	s := string(param)
	switch s {
	case "1":
		return Key{Type: KeyHome}, nil
	case "3":
		return Key{Type: KeyDelete}, nil
	case "4":
		return Key{Type: KeyEnd}, nil
	case "5":
		return Key{Type: KeyPageUp}, nil
	case "6":
		return Key{Type: KeyPageDown}, nil
	default:
		return Key{Type: KeyEscape}, nil
	}
}

func (r *Reader) parseUTF8(first byte) (Key, error) {
	// Determine how many continuation bytes we need
	var size int
	switch {
	case first&0xe0 == 0xc0:
		size = 2
	case first&0xf0 == 0xe0:
		size = 3
	case first&0xf8 == 0xf0:
		size = 4
	default:
		return Key{Type: KeyRune, Rune: rune(first)}, nil
	}

	buf := make([]byte, size)
	buf[0] = first
	for i := 1; i < size; i++ {
		b, err := r.readByte()
		if err != nil {
			return Key{Type: KeyRune, Rune: rune(first)}, nil
		}
		buf[i] = b
	}

	r2 := rune(0)
	switch size {
	case 2:
		r2 = rune(buf[0]&0x1f)<<6 | rune(buf[1]&0x3f)
	case 3:
		r2 = rune(buf[0]&0x0f)<<12 | rune(buf[1]&0x3f)<<6 | rune(buf[2]&0x3f)
	case 4:
		r2 = rune(buf[0]&0x07)<<18 | rune(buf[1]&0x3f)<<12 | rune(buf[2]&0x3f)<<6 | rune(buf[3]&0x3f)
	}

	return Key{Type: KeyRune, Rune: r2}, nil
}

func (r *Reader) readByte() (byte, error) {
	if r.pos < r.n {
		b := r.buf[r.pos]
		r.pos++
		return b, nil
	}
	n, err := r.in.Read(r.buf[:])
	if err != nil {
		return 0, err
	}
	r.n = n
	r.pos = 1
	return r.buf[0], nil
}

func (r *Reader) peekByte() (byte, error) {
	if r.pos < r.n {
		return r.buf[r.pos], nil
	}
	n, err := r.in.Read(r.buf[:])
	if err != nil {
		return 0, err
	}
	r.n = n
	r.pos = 0
	return r.buf[0], nil
}

func (r *Reader) consumeByte() {
	r.pos++
}
