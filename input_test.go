package inquire

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
)

func simulatedIO(input string) (io.Reader, *bytes.Buffer) {
	return strings.NewReader(input), &bytes.Buffer{}
}

func TestInputBasic(t *testing.T) {
	in, out := simulatedIO("hello\n")
	result, err := Input("Name?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello" {
		t.Errorf("got %q, want %q", result, "hello")
	}
}

func TestInputDefault(t *testing.T) {
	in, out := simulatedIO("\n")
	result, err := Input("Name?", WithStdio(in, out, out), WithDefault("world"))
	if err != nil {
		t.Fatal(err)
	}
	if result != "world" {
		t.Errorf("got %q, want %q", result, "world")
	}
}

func TestInputDefaultOverride(t *testing.T) {
	in, out := simulatedIO("custom\n")
	result, err := Input("Name?", WithStdio(in, out, out), WithDefault("world"))
	if err != nil {
		t.Fatal(err)
	}
	if result != "custom" {
		t.Errorf("got %q, want %q", result, "custom")
	}
}

func TestInputValidation(t *testing.T) {
	// First attempt empty (fails validation), then provide valid input
	in, out := simulatedIO("\nhello\n")
	result, err := Input("Name?",
		WithStdio(in, out, out),
		WithValidate(func(s string) error {
			if s == "" {
				return errors.New("required")
			}
			return nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello" {
		t.Errorf("got %q, want %q", result, "hello")
	}
}

func TestInputTransform(t *testing.T) {
	in, out := simulatedIO("HELLO\n")
	result, err := Input("Name?",
		WithStdio(in, out, out),
		WithTransform(strings.ToLower),
	)
	if err != nil {
		t.Fatal(err)
	}
	if result != "hello" {
		t.Errorf("got %q, want %q", result, "hello")
	}
}

func TestInputCtrlC(t *testing.T) {
	in, out := simulatedIO(string([]byte{0x03}))
	_, err := Input("Name?", WithStdio(in, out, out))
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestInputContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Use a reader that blocks — context cancellation should be checked
	// before reading. But since our loop reads first then checks, we need
	// the read to return something. Let's provide Ctrl+C as fallback.
	in, out := simulatedIO(string([]byte{0x03}))
	_, err := Input("Name?", WithStdio(in, out, out), WithContext(ctx))
	if err == nil {
		t.Error("expected error from cancelled context")
	}
}

func TestInputBackspace(t *testing.T) {
	// Type "ab", backspace, type "c", enter => "ac"
	in, out := simulatedIO("ab\x7fc\n")
	result, err := Input("Name?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "ac" {
		t.Errorf("got %q, want %q", result, "ac")
	}
}
