package promptui

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }

func stdinFrom(s string) io.ReadCloser {
	return nopCloser{strings.NewReader(s)}
}

func stdoutDiscard() io.WriteCloser {
	return nopWriteCloser{io.Discard}
}

func TestPromptInput(t *testing.T) {
	p := &Prompt{
		Label:  "Name",
		Stdin:  stdinFrom("Alice\n"),
		Stdout: stdoutDiscard(),
	}
	result, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != "Alice" {
		t.Errorf("got %q, want %q", result, "Alice")
	}
}

func TestPromptDefault(t *testing.T) {
	p := &Prompt{
		Label:   "Name",
		Default: "Bob",
		Stdin:   stdinFrom("\n"),
		Stdout:  stdoutDiscard(),
	}
	result, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != "Bob" {
		t.Errorf("got %q, want %q", result, "Bob")
	}
}

func TestPromptValidation(t *testing.T) {
	// Empty input fails, then valid input
	p := &Prompt{
		Label: "Name",
		Validate: func(s string) error {
			if s == "" {
				return errors.New("required")
			}
			return nil
		},
		Stdin:  stdinFrom("\nAlice\n"),
		Stdout: stdoutDiscard(),
	}
	result, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != "Alice" {
		t.Errorf("got %q, want %q", result, "Alice")
	}
}

func TestPromptMask(t *testing.T) {
	p := &Prompt{
		Label:  "Password",
		Mask:   '*',
		Stdin:  stdinFrom("secret\n"),
		Stdout: stdoutDiscard(),
	}
	result, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != "secret" {
		t.Errorf("got %q, want %q", result, "secret")
	}
}

func TestPromptConfirm(t *testing.T) {
	p := &Prompt{
		Label:     "Continue?",
		IsConfirm: true,
		Stdin:     stdinFrom("y\n"),
		Stdout:    stdoutDiscard(),
	}
	result, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != "y" {
		t.Errorf("got %q, want %q", result, "y")
	}
}

func TestPromptConfirmNo(t *testing.T) {
	p := &Prompt{
		Label:     "Continue?",
		IsConfirm: true,
		Stdin:     stdinFrom("n\n"),
		Stdout:    stdoutDiscard(),
	}
	result, err := p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if result != "N" {
		t.Errorf("got %q, want %q", result, "N")
	}
}

func TestPromptCtrlC(t *testing.T) {
	p := &Prompt{
		Label:  "Name",
		Stdin:  stdinFrom(string([]byte{0x03})),
		Stdout: stdoutDiscard(),
	}
	_, err := p.Run()
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestSelectBasic(t *testing.T) {
	// Press Enter to select first
	s := &Select{
		Label:  "Color",
		Items:  []string{"red", "green", "blue"},
		Stdin:  stdinFrom("\n"),
		Stdout: stdoutDiscard(),
	}
	idx, val, err := s.Run()
	if err != nil {
		t.Fatal(err)
	}
	if idx != 0 || val != "red" {
		t.Errorf("got (%d, %q), want (0, red)", idx, val)
	}
}

func TestSelectWithSize(t *testing.T) {
	s := &Select{
		Label:  "Color",
		Items:  []string{"red", "green", "blue"},
		Size:   2,
		Stdin:  stdinFrom("\n"),
		Stdout: stdoutDiscard(),
	}
	_, val, err := s.Run()
	if err != nil {
		t.Fatal(err)
	}
	if val != "red" {
		t.Errorf("got %q, want %q", val, "red")
	}
}

func TestSelectCtrlC(t *testing.T) {
	s := &Select{
		Label:  "Color",
		Items:  []string{"red"},
		Stdin:  stdinFrom(string([]byte{0x03})),
		Stdout: stdoutDiscard(),
	}
	_, _, err := s.Run()
	if !errors.Is(err, ErrInterrupt) {
		t.Errorf("got %v, want ErrInterrupt", err)
	}
}

func TestSelectTemplates(t *testing.T) {
	// Templates are accepted but ignored in compat mode
	s := &Select{
		Label: "Color",
		Items: []string{"red", "green"},
		Templates: &SelectTemplates{
			Active:   "> {{ . }}",
			Inactive: "  {{ . }}",
		},
		Stdin:  stdinFrom("\n"),
		Stdout: stdoutDiscard(),
	}
	_, val, err := s.Run()
	if err != nil {
		t.Fatal(err)
	}
	if val != "red" {
		t.Errorf("got %q, want %q", val, "red")
	}
}

func TestErrInterrupt(t *testing.T) {
	if !errors.Is(ErrInterrupt, ErrInterrupt) {
		t.Error("ErrInterrupt should match itself")
	}
}
