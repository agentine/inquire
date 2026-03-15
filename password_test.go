package inquire

import (
	"errors"
	"strings"
	"testing"
)

func TestPasswordBasic(t *testing.T) {
	in, out := simulatedIO("secret\n")
	result, err := Password("Password?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "secret" {
		t.Errorf("got %q, want %q", result, "secret")
	}
}

func TestPasswordMasked(t *testing.T) {
	in, out := simulatedIO("abc\n")
	result, err := Password("Password?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "abc" {
		t.Errorf("got %q, want %q", result, "abc")
	}
	// Output should contain asterisks, not the actual password
	output := out.String()
	if strings.Contains(output, "abc") && !strings.Contains(output, "***") {
		t.Error("output should show asterisks, not plaintext")
	}
}

func TestPasswordCtrlC(t *testing.T) {
	in, out := simulatedIO(string([]byte{0x03}))
	_, err := Password("Password?", WithStdio(in, out, out))
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestPasswordValidation(t *testing.T) {
	// Empty password fails, then provide valid one
	in, out := simulatedIO("\nsecret\n")
	result, err := Password("Password?",
		WithStdio(in, out, out),
		WithValidate(func(s string) error {
			if len(s) < 3 {
				return errors.New("too short")
			}
			return nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result != "secret" {
		t.Errorf("got %q, want %q", result, "secret")
	}
}

func TestPasswordBackspace(t *testing.T) {
	// Type "ab", backspace, type "c", enter => "ac"
	in, out := simulatedIO("ab\x7fc\n")
	result, err := Password("Password?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "ac" {
		t.Errorf("got %q, want %q", result, "ac")
	}
}
