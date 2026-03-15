package inquire

import (
	"errors"
	"testing"
)

func TestSelectBasic(t *testing.T) {
	// Down arrow to index 1, then Enter
	in, out := simulatedIO(string([]byte{0x1b, '[', 'B', 0x0d}))
	idx, val, err := Select("Color?", []string{"red", "green", "blue"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if idx != 1 || val != "green" {
		t.Errorf("got (%d, %q), want (1, green)", idx, val)
	}
}

func TestSelectFirst(t *testing.T) {
	// Just press Enter to select first item
	in, out := simulatedIO("\n")
	idx, val, err := Select("Color?", []string{"red", "green", "blue"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if idx != 0 || val != "red" {
		t.Errorf("got (%d, %q), want (0, red)", idx, val)
	}
}

func TestSelectLast(t *testing.T) {
	// Down, Down, Enter
	in, out := simulatedIO(string([]byte{0x1b, '[', 'B', 0x1b, '[', 'B', 0x0d}))
	idx, val, err := Select("Color?", []string{"red", "green", "blue"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if idx != 2 || val != "blue" {
		t.Errorf("got (%d, %q), want (2, blue)", idx, val)
	}
}

func TestSelectDefault(t *testing.T) {
	// Press Enter with default set to index 2
	in, out := simulatedIO("\n")
	idx, val, err := Select("Color?", []string{"red", "green", "blue"},
		WithStdio(in, out, out), WithDefault(2))
	if err != nil {
		t.Fatal(err)
	}
	if idx != 2 || val != "blue" {
		t.Errorf("got (%d, %q), want (2, blue)", idx, val)
	}
}

func TestSelectDefaultString(t *testing.T) {
	in, out := simulatedIO("\n")
	idx, val, err := Select("Color?", []string{"red", "green", "blue"},
		WithStdio(in, out, out), WithDefault("green"))
	if err != nil {
		t.Fatal(err)
	}
	if idx != 1 || val != "green" {
		t.Errorf("got (%d, %q), want (1, green)", idx, val)
	}
}

func TestSelectCtrlC(t *testing.T) {
	in, out := simulatedIO(string([]byte{0x03}))
	_, _, err := Select("Color?", []string{"red", "green"}, WithStdio(in, out, out))
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestSelectEmpty(t *testing.T) {
	in, out := simulatedIO("\n")
	_, _, err := Select("Color?", []string{}, WithStdio(in, out, out))
	if err == nil {
		t.Error("expected error for empty options")
	}
}

func TestSelectFilter(t *testing.T) {
	// Type "gr" to filter to "green", then Enter
	in, out := simulatedIO("gr\n")
	idx, val, err := Select("Color?", []string{"red", "green", "blue"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	// idx should be the original index of "green"
	if idx != 1 || val != "green" {
		t.Errorf("got (%d, %q), want (1, green)", idx, val)
	}
}
