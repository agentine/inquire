// Package term provides raw terminal mode handling and ANSI escape helpers.
package term

import (
	"fmt"
	"io"
	"os"
)

// State holds the saved terminal state for restoring later.
type State struct {
	fd    int
	saved interface{} // platform-specific saved state
}

// MakeRaw puts the terminal attached to fd into raw mode and returns
// the saved state for later restoration via Restore.
func MakeRaw(fd int) (*State, error) {
	return makeRaw(fd)
}

// Restore restores the terminal to a previously saved state.
func Restore(s *State) error {
	return restore(s)
}

// IsTerminal reports whether fd is a terminal.
func IsTerminal(fd int) bool {
	return isTerminal(fd)
}

// Size returns the dimensions of the terminal attached to fd.
func Size(fd int) (width, height int, err error) {
	return size(fd)
}

// ANSI escape sequence helpers.

// CursorUp moves the cursor up n lines.
func CursorUp(w io.Writer, n int) {
	if n > 0 {
		fmt.Fprintf(w, "\x1b[%dA", n)
	}
}

// CursorDown moves the cursor down n lines.
func CursorDown(w io.Writer, n int) {
	if n > 0 {
		fmt.Fprintf(w, "\x1b[%dB", n)
	}
}

// CursorForward moves the cursor right n columns.
func CursorForward(w io.Writer, n int) {
	if n > 0 {
		fmt.Fprintf(w, "\x1b[%dC", n)
	}
}

// CursorBack moves the cursor left n columns.
func CursorBack(w io.Writer, n int) {
	if n > 0 {
		fmt.Fprintf(w, "\x1b[%dD", n)
	}
}

// CursorColumn moves the cursor to column n (1-based).
func CursorColumn(w io.Writer, n int) {
	fmt.Fprintf(w, "\x1b[%dG", n)
}

// ClearLine clears the current line.
func ClearLine(w io.Writer) {
	fmt.Fprint(w, "\x1b[2K")
}

// ClearToEnd clears from cursor to end of line.
func ClearToEnd(w io.Writer) {
	fmt.Fprint(w, "\x1b[K")
}

// ClearScreen clears the entire screen and moves cursor to top-left.
func ClearScreen(w io.Writer) {
	fmt.Fprint(w, "\x1b[2J\x1b[H")
}

// HideCursor hides the terminal cursor.
func HideCursor(w io.Writer) {
	fmt.Fprint(w, "\x1b[?25l")
}

// ShowCursor shows the terminal cursor.
func ShowCursor(w io.Writer) {
	fmt.Fprint(w, "\x1b[?25h")
}

// SaveCursor saves the cursor position.
func SaveCursor(w io.Writer) {
	fmt.Fprint(w, "\x1b7")
}

// RestoreCursor restores the cursor position.
func RestoreCursor(w io.Writer) {
	fmt.Fprint(w, "\x1b8")
}

// Style constants for ANSI SGR.
const (
	Reset     = "\x1b[0m"
	Bold      = "\x1b[1m"
	Dim       = "\x1b[2m"
	Italic    = "\x1b[3m"
	Underline = "\x1b[4m"

	FgBlack   = "\x1b[30m"
	FgRed     = "\x1b[31m"
	FgGreen   = "\x1b[32m"
	FgYellow  = "\x1b[33m"
	FgBlue    = "\x1b[34m"
	FgMagenta = "\x1b[35m"
	FgCyan    = "\x1b[36m"
	FgWhite   = "\x1b[37m"

	FgDefault = "\x1b[39m"
)

// Stdio holds overridable I/O streams for prompts.
type Stdio struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

// DefaultStdio returns the standard I/O streams.
func DefaultStdio() Stdio {
	return Stdio{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}
