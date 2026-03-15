package inquire

import (
	"errors"
	"testing"
)

func TestConfirmYes(t *testing.T) {
	in, out := simulatedIO("y\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("expected true for 'y'")
	}
}

func TestConfirmYesFull(t *testing.T) {
	in, out := simulatedIO("yes\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("expected true for 'yes'")
	}
}

func TestConfirmNo(t *testing.T) {
	in, out := simulatedIO("n\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("expected false for 'n'")
	}
}

func TestConfirmNoFull(t *testing.T) {
	in, out := simulatedIO("no\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("expected false for 'no'")
	}
}

func TestConfirmDefaultFalse(t *testing.T) {
	// Just press enter with no default — should be false
	in, out := simulatedIO("\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("expected false for empty input with no default")
	}
}

func TestConfirmDefaultTrue(t *testing.T) {
	in, out := simulatedIO("\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out), WithDefault(true))
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("expected true for empty input with default=true")
	}
}

func TestConfirmDefaultOverride(t *testing.T) {
	in, out := simulatedIO("n\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out), WithDefault(true))
	if err != nil {
		t.Fatal(err)
	}
	if result {
		t.Error("expected false for 'n' even with default=true")
	}
}

func TestConfirmCtrlC(t *testing.T) {
	in, out := simulatedIO(string([]byte{0x03}))
	_, err := Confirm("Continue?", WithStdio(in, out, out))
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestConfirmCaseInsensitive(t *testing.T) {
	in, out := simulatedIO("Y\n")
	result, err := Confirm("Continue?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Error("expected true for 'Y'")
	}
}
