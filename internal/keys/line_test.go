package keys

import "testing"

func TestLineInsert(t *testing.T) {
	l := NewLine()
	l.Insert('h')
	l.Insert('i')
	if got := l.String(); got != "hi" {
		t.Errorf("got %q, want %q", got, "hi")
	}
	if l.Pos() != 2 {
		t.Errorf("pos = %d, want 2", l.Pos())
	}
}

func TestLineInsertMiddle(t *testing.T) {
	l := NewLineFrom("ac")
	l.MoveLeft()
	l.Insert('b')
	if got := l.String(); got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
	if l.Pos() != 2 {
		t.Errorf("pos = %d, want 2", l.Pos())
	}
}

func TestLineBackspace(t *testing.T) {
	l := NewLineFrom("abc")
	ok := l.Backspace()
	if !ok {
		t.Error("Backspace returned false")
	}
	if got := l.String(); got != "ab" {
		t.Errorf("got %q, want %q", got, "ab")
	}
}

func TestLineBackspaceEmpty(t *testing.T) {
	l := NewLine()
	if l.Backspace() {
		t.Error("Backspace on empty returned true")
	}
}

func TestLineDelete(t *testing.T) {
	l := NewLineFrom("abc")
	l.MoveHome()
	ok := l.Delete()
	if !ok {
		t.Error("Delete returned false")
	}
	if got := l.String(); got != "bc" {
		t.Errorf("got %q, want %q", got, "bc")
	}
}

func TestLineDeleteAtEnd(t *testing.T) {
	l := NewLineFrom("abc")
	if l.Delete() {
		t.Error("Delete at end returned true")
	}
}

func TestLineMovement(t *testing.T) {
	l := NewLineFrom("abc")
	if l.Pos() != 3 {
		t.Errorf("initial pos = %d, want 3", l.Pos())
	}
	l.MoveLeft()
	if l.Pos() != 2 {
		t.Errorf("after left pos = %d, want 2", l.Pos())
	}
	l.MoveRight()
	if l.Pos() != 3 {
		t.Errorf("after right pos = %d, want 3", l.Pos())
	}
	l.MoveHome()
	if l.Pos() != 0 {
		t.Errorf("after home pos = %d, want 0", l.Pos())
	}
	l.MoveEnd()
	if l.Pos() != 3 {
		t.Errorf("after end pos = %d, want 3", l.Pos())
	}
}

func TestLineMoveLeftBoundary(t *testing.T) {
	l := NewLine()
	if l.MoveLeft() {
		t.Error("MoveLeft on empty returned true")
	}
}

func TestLineMoveRightBoundary(t *testing.T) {
	l := NewLine()
	if l.MoveRight() {
		t.Error("MoveRight on empty returned true")
	}
}

func TestLineDeleteToEnd(t *testing.T) {
	l := NewLineFrom("hello world")
	l.MoveHome()
	l.MoveRight()
	l.MoveRight()
	l.MoveRight()
	l.MoveRight()
	l.MoveRight()
	l.DeleteToEnd()
	if got := l.String(); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestLineDeleteToStart(t *testing.T) {
	l := NewLineFrom("hello world")
	// Move cursor to position 6 ("hello " consumed)
	l.MoveHome()
	for i := 0; i < 6; i++ {
		l.MoveRight()
	}
	l.DeleteToStart()
	if got := l.String(); got != "world" {
		t.Errorf("got %q, want %q", got, "world")
	}
	if l.Pos() != 0 {
		t.Errorf("pos = %d, want 0", l.Pos())
	}
}

func TestLineDeleteWord(t *testing.T) {
	l := NewLineFrom("hello world")
	ok := l.DeleteWord()
	if !ok {
		t.Error("DeleteWord returned false")
	}
	if got := l.String(); got != "hello " {
		t.Errorf("got %q, want %q", got, "hello ")
	}
}

func TestLineDeleteWordEmpty(t *testing.T) {
	l := NewLine()
	if l.DeleteWord() {
		t.Error("DeleteWord on empty returned true")
	}
}

func TestLineSet(t *testing.T) {
	l := NewLine()
	l.Set("new content")
	if got := l.String(); got != "new content" {
		t.Errorf("got %q, want %q", got, "new content")
	}
	if l.Pos() != 11 {
		t.Errorf("pos = %d, want 11", l.Pos())
	}
}

func TestLineClear(t *testing.T) {
	l := NewLineFrom("hello")
	l.Clear()
	if got := l.String(); got != "" {
		t.Errorf("got %q, want empty", got)
	}
	if l.Pos() != 0 {
		t.Errorf("pos = %d, want 0", l.Pos())
	}
}

func TestLineLen(t *testing.T) {
	l := NewLineFrom("hi")
	if l.Len() != 2 {
		t.Errorf("Len() = %d, want 2", l.Len())
	}
}

func TestLineUnicode(t *testing.T) {
	l := NewLine()
	l.Insert('é')
	l.Insert('ñ')
	if got := l.String(); got != "éñ" {
		t.Errorf("got %q, want %q", got, "éñ")
	}
	if l.Len() != 2 {
		t.Errorf("Len() = %d, want 2", l.Len())
	}
}
