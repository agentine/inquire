package keys

import (
	"bytes"
	"testing"
)

func TestReadKeyPrintable(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte("abc")))
	for _, want := range []rune{'a', 'b', 'c'} {
		k, err := r.ReadKey()
		if err != nil {
			t.Fatal(err)
		}
		if k.Type != KeyRune || k.Rune != want {
			t.Errorf("got %v %c, want Rune %c", k.Type, k.Rune, want)
		}
	}
}

func TestReadKeyEnter(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x0d}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyEnter {
		t.Errorf("got %v, want Enter", k.Type)
	}
}

func TestReadKeyLineFeed(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x0a}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyEnter {
		t.Errorf("got %v, want Enter", k.Type)
	}
}

func TestReadKeyBackspace(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x7f}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyBackspace {
		t.Errorf("got %v, want Backspace", k.Type)
	}
}

func TestReadKeyTab(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x09}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyTab {
		t.Errorf("got %v, want Tab", k.Type)
	}
}

func TestReadKeyCtrlC(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x03}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyCtrlC {
		t.Errorf("got %v, want Ctrl+C", k.Type)
	}
}

func TestReadKeyCtrlD(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x04}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyCtrlD {
		t.Errorf("got %v, want Ctrl+D", k.Type)
	}
}

func TestReadKeyArrows(t *testing.T) {
	tests := []struct {
		input []byte
		want  KeyType
	}{
		{[]byte{0x1b, '[', 'A'}, KeyUp},
		{[]byte{0x1b, '[', 'B'}, KeyDown},
		{[]byte{0x1b, '[', 'C'}, KeyRight},
		{[]byte{0x1b, '[', 'D'}, KeyLeft},
	}
	for _, tt := range tests {
		r := NewReader(bytes.NewReader(tt.input))
		k, err := r.ReadKey()
		if err != nil {
			t.Fatal(err)
		}
		if k.Type != tt.want {
			t.Errorf("input %v: got %v, want %v", tt.input, k.Type, tt.want)
		}
	}
}

func TestReadKeyHomeEnd(t *testing.T) {
	tests := []struct {
		input []byte
		want  KeyType
	}{
		{[]byte{0x1b, '[', 'H'}, KeyHome},
		{[]byte{0x1b, '[', 'F'}, KeyEnd},
		{[]byte{0x1b, 'O', 'H'}, KeyHome},
		{[]byte{0x1b, 'O', 'F'}, KeyEnd},
		{[]byte{0x1b, '[', '1', '~'}, KeyHome},
		{[]byte{0x1b, '[', '4', '~'}, KeyEnd},
	}
	for _, tt := range tests {
		r := NewReader(bytes.NewReader(tt.input))
		k, err := r.ReadKey()
		if err != nil {
			t.Fatal(err)
		}
		if k.Type != tt.want {
			t.Errorf("input %v: got %v, want %v", tt.input, k.Type, tt.want)
		}
	}
}

func TestReadKeyDelete(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{0x1b, '[', '3', '~'}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyDelete {
		t.Errorf("got %v, want Delete", k.Type)
	}
}

func TestReadKeyPageUpDown(t *testing.T) {
	tests := []struct {
		input []byte
		want  KeyType
	}{
		{[]byte{0x1b, '[', '5', '~'}, KeyPageUp},
		{[]byte{0x1b, '[', '6', '~'}, KeyPageDown},
	}
	for _, tt := range tests {
		r := NewReader(bytes.NewReader(tt.input))
		k, err := r.ReadKey()
		if err != nil {
			t.Fatal(err)
		}
		if k.Type != tt.want {
			t.Errorf("input %v: got %v, want %v", tt.input, k.Type, tt.want)
		}
	}
}

func TestReadKeySpace(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte{' '}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeySpace {
		t.Errorf("got %v, want Space", k.Type)
	}
	if k.Rune != ' ' {
		t.Errorf("got rune %c, want space", k.Rune)
	}
}

func TestReadKeyUTF8(t *testing.T) {
	// "é" is 0xc3 0xa9 in UTF-8
	r := NewReader(bytes.NewReader([]byte{0xc3, 0xa9}))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyRune || k.Rune != 'é' {
		t.Errorf("got %v %c, want Rune é", k.Type, k.Rune)
	}
}

func TestReadKeyUTF8Emoji(t *testing.T) {
	// "😀" is 4 bytes in UTF-8
	input := []byte(string([]rune{'😀'}))
	r := NewReader(bytes.NewReader(input))
	k, err := r.ReadKey()
	if err != nil {
		t.Fatal(err)
	}
	if k.Type != KeyRune || k.Rune != '😀' {
		t.Errorf("got %v %U, want Rune U+1F600", k.Type, k.Rune)
	}
}

func TestReadKeyCtrlSequences(t *testing.T) {
	tests := []struct {
		input byte
		want  KeyType
	}{
		{0x01, KeyCtrlA},
		{0x02, KeyCtrlB},
		{0x05, KeyCtrlE},
		{0x06, KeyCtrlF},
		{0x0b, KeyCtrlK},
		{0x0c, KeyCtrlL},
		{0x0e, KeyCtrlN},
		{0x10, KeyCtrlP},
		{0x15, KeyCtrlU},
		{0x17, KeyCtrlW},
	}
	for _, tt := range tests {
		r := NewReader(bytes.NewReader([]byte{tt.input}))
		k, err := r.ReadKey()
		if err != nil {
			t.Fatal(err)
		}
		if k.Type != tt.want {
			t.Errorf("input 0x%02x: got %v, want %v", tt.input, k.Type, tt.want)
		}
	}
}

func TestKeyTypeString(t *testing.T) {
	tests := []struct {
		k    KeyType
		want string
	}{
		{KeyRune, "Rune"},
		{KeyEnter, "Enter"},
		{KeyBackspace, "Backspace"},
		{KeyUp, "Up"},
		{KeyCtrlC, "Ctrl+C"},
	}
	for _, tt := range tests {
		if got := tt.k.String(); got != tt.want {
			t.Errorf("%d.String() = %q, want %q", tt.k, got, tt.want)
		}
	}
}
