package inquire

import (
	"errors"
	"testing"
)

func TestMultiSelectBasic(t *testing.T) {
	// Space to toggle first, Down, Space to toggle second, Enter
	in, out := simulatedIO(string([]byte{' ', 0x1b, '[', 'B', ' ', 0x0d}))
	indices, vals, err := MultiSelect("Colors?", []string{"red", "green", "blue"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if len(indices) != 2 || indices[0] != 0 || indices[1] != 1 {
		t.Errorf("indices = %v, want [0 1]", indices)
	}
	if len(vals) != 2 || vals[0] != "red" || vals[1] != "green" {
		t.Errorf("vals = %v, want [red green]", vals)
	}
}

func TestMultiSelectNone(t *testing.T) {
	// Just press Enter without selecting anything
	in, out := simulatedIO("\n")
	indices, vals, err := MultiSelect("Colors?", []string{"red", "green"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if len(indices) != 0 {
		t.Errorf("indices = %v, want empty", indices)
	}
	if len(vals) != 0 {
		t.Errorf("vals = %v, want empty", vals)
	}
}

func TestMultiSelectToggle(t *testing.T) {
	// Select first, then deselect it, then Enter — should be empty
	in, out := simulatedIO(string([]byte{' ', ' ', 0x0d}))
	indices, vals, err := MultiSelect("Colors?", []string{"red", "green"}, WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if len(indices) != 0 {
		t.Errorf("indices = %v, want empty", indices)
	}
	if len(vals) != 0 {
		t.Errorf("vals = %v, want empty", vals)
	}
}

func TestMultiSelectCtrlC(t *testing.T) {
	in, out := simulatedIO(string([]byte{0x03}))
	_, _, err := MultiSelect("Colors?", []string{"red"}, WithStdio(in, out, out))
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestMultiSelectEmpty(t *testing.T) {
	in, out := simulatedIO("\n")
	_, _, err := MultiSelect("Colors?", []string{}, WithStdio(in, out, out))
	if err == nil {
		t.Error("expected error for empty options")
	}
}

func TestMultiSelectDefaults(t *testing.T) {
	// Just press Enter with defaults pre-selected
	in, out := simulatedIO("\n")
	indices, vals, err := MultiSelect("Colors?", []string{"red", "green", "blue"},
		WithStdio(in, out, out), WithDefault([]int{0, 2}))
	if err != nil {
		t.Fatal(err)
	}
	if len(indices) != 2 || indices[0] != 0 || indices[1] != 2 {
		t.Errorf("indices = %v, want [0 2]", indices)
	}
	if len(vals) != 2 || vals[0] != "red" || vals[1] != "blue" {
		t.Errorf("vals = %v, want [red blue]", vals)
	}
}
