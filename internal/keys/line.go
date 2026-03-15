package keys

import "unicode/utf8"

// Line provides basic line editing (insert, delete, cursor movement).
type Line struct {
	buf []rune
	pos int // cursor position (rune index)
}

// NewLine creates an empty line editor.
func NewLine() *Line {
	return &Line{}
}

// NewLineFrom creates a line editor with initial content.
func NewLineFrom(s string) *Line {
	runes := []rune(s)
	return &Line{buf: runes, pos: len(runes)}
}

// Insert inserts a rune at the cursor position.
func (l *Line) Insert(r rune) {
	if !utf8.ValidRune(r) {
		return
	}
	if l.pos == len(l.buf) {
		l.buf = append(l.buf, r)
	} else {
		l.buf = append(l.buf, 0)
		copy(l.buf[l.pos+1:], l.buf[l.pos:])
		l.buf[l.pos] = r
	}
	l.pos++
}

// Backspace deletes the rune before the cursor.
func (l *Line) Backspace() bool {
	if l.pos == 0 {
		return false
	}
	l.pos--
	l.buf = append(l.buf[:l.pos], l.buf[l.pos+1:]...)
	return true
}

// Delete deletes the rune at the cursor.
func (l *Line) Delete() bool {
	if l.pos >= len(l.buf) {
		return false
	}
	l.buf = append(l.buf[:l.pos], l.buf[l.pos+1:]...)
	return true
}

// MoveLeft moves the cursor one position left.
func (l *Line) MoveLeft() bool {
	if l.pos > 0 {
		l.pos--
		return true
	}
	return false
}

// MoveRight moves the cursor one position right.
func (l *Line) MoveRight() bool {
	if l.pos < len(l.buf) {
		l.pos++
		return true
	}
	return false
}

// MoveHome moves the cursor to the beginning.
func (l *Line) MoveHome() {
	l.pos = 0
}

// MoveEnd moves the cursor to the end.
func (l *Line) MoveEnd() {
	l.pos = len(l.buf)
}

// DeleteToEnd deletes from cursor to end of line.
func (l *Line) DeleteToEnd() {
	l.buf = l.buf[:l.pos]
}

// DeleteToStart deletes from start to cursor.
func (l *Line) DeleteToStart() {
	l.buf = l.buf[l.pos:]
	l.pos = 0
}

// DeleteWord deletes the word before the cursor.
func (l *Line) DeleteWord() bool {
	if l.pos == 0 {
		return false
	}
	// Skip trailing spaces
	end := l.pos
	for l.pos > 0 && l.buf[l.pos-1] == ' ' {
		l.pos--
	}
	// Delete until next space
	for l.pos > 0 && l.buf[l.pos-1] != ' ' {
		l.pos--
	}
	l.buf = append(l.buf[:l.pos], l.buf[end:]...)
	return true
}

// String returns the current line content.
func (l *Line) String() string {
	return string(l.buf)
}

// Pos returns the cursor position.
func (l *Line) Pos() int {
	return l.pos
}

// Len returns the number of runes in the line.
func (l *Line) Len() int {
	return len(l.buf)
}

// Set replaces the line content and moves cursor to end.
func (l *Line) Set(s string) {
	l.buf = []rune(s)
	l.pos = len(l.buf)
}

// Clear empties the line.
func (l *Line) Clear() {
	l.buf = l.buf[:0]
	l.pos = 0
}
