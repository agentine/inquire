package term

import (
	"bytes"
	"testing"
)

func TestCursorUp(t *testing.T) {
	var buf bytes.Buffer
	CursorUp(&buf, 3)
	if got := buf.String(); got != "\x1b[3A" {
		t.Errorf("got %q, want %q", got, "\x1b[3A")
	}
}

func TestCursorUpZero(t *testing.T) {
	var buf bytes.Buffer
	CursorUp(&buf, 0)
	if buf.Len() != 0 {
		t.Error("CursorUp(0) should produce no output")
	}
}

func TestCursorDown(t *testing.T) {
	var buf bytes.Buffer
	CursorDown(&buf, 2)
	if got := buf.String(); got != "\x1b[2B" {
		t.Errorf("got %q, want %q", got, "\x1b[2B")
	}
}

func TestCursorForward(t *testing.T) {
	var buf bytes.Buffer
	CursorForward(&buf, 5)
	if got := buf.String(); got != "\x1b[5C" {
		t.Errorf("got %q, want %q", got, "\x1b[5C")
	}
}

func TestCursorBack(t *testing.T) {
	var buf bytes.Buffer
	CursorBack(&buf, 1)
	if got := buf.String(); got != "\x1b[1D" {
		t.Errorf("got %q, want %q", got, "\x1b[1D")
	}
}

func TestCursorColumn(t *testing.T) {
	var buf bytes.Buffer
	CursorColumn(&buf, 10)
	if got := buf.String(); got != "\x1b[10G" {
		t.Errorf("got %q, want %q", got, "\x1b[10G")
	}
}

func TestClearLine(t *testing.T) {
	var buf bytes.Buffer
	ClearLine(&buf)
	if got := buf.String(); got != "\x1b[2K" {
		t.Errorf("got %q, want %q", got, "\x1b[2K")
	}
}

func TestClearToEnd(t *testing.T) {
	var buf bytes.Buffer
	ClearToEnd(&buf)
	if got := buf.String(); got != "\x1b[K" {
		t.Errorf("got %q, want %q", got, "\x1b[K")
	}
}

func TestClearScreen(t *testing.T) {
	var buf bytes.Buffer
	ClearScreen(&buf)
	if got := buf.String(); got != "\x1b[2J\x1b[H" {
		t.Errorf("got %q, want %q", got, "\x1b[2J\x1b[H")
	}
}

func TestHideShowCursor(t *testing.T) {
	var buf bytes.Buffer
	HideCursor(&buf)
	if got := buf.String(); got != "\x1b[?25l" {
		t.Errorf("HideCursor got %q", got)
	}
	buf.Reset()
	ShowCursor(&buf)
	if got := buf.String(); got != "\x1b[?25h" {
		t.Errorf("ShowCursor got %q", got)
	}
}

func TestSaveRestoreCursor(t *testing.T) {
	var buf bytes.Buffer
	SaveCursor(&buf)
	if got := buf.String(); got != "\x1b7" {
		t.Errorf("SaveCursor got %q", got)
	}
	buf.Reset()
	RestoreCursor(&buf)
	if got := buf.String(); got != "\x1b8" {
		t.Errorf("RestoreCursor got %q", got)
	}
}

func TestDefaultStdio(t *testing.T) {
	s := DefaultStdio()
	if s.In == nil || s.Out == nil || s.Err == nil {
		t.Error("DefaultStdio returned nil streams")
	}
}

func TestStyleConstants(t *testing.T) {
	// Verify they're non-empty escape sequences
	styles := []string{Reset, Bold, Dim, Italic, Underline, FgRed, FgGreen, FgBlue, FgDefault}
	for _, s := range styles {
		if len(s) == 0 {
			t.Error("style constant is empty")
		}
		if s[0] != 0x1b {
			t.Errorf("style %q doesn't start with ESC", s)
		}
	}
}
