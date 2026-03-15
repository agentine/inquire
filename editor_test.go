package inquire

import (
	"strings"
	"testing"
)

func TestEditorSimulated(t *testing.T) {
	in, out := simulatedIO("some editor content\n")
	result, err := Editor("Description?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "some editor content" {
		t.Errorf("got %q, want %q", result, "some editor content")
	}
}

func TestEditorSimulatedMultiline(t *testing.T) {
	in, out := simulatedIO("line1\nline2\nline3\n")
	result, err := Editor("Description?", WithStdio(in, out, out))
	if err != nil {
		t.Fatal(err)
	}
	if result != "line1\nline2\nline3" {
		t.Errorf("got %q, want %q", result, "line1\nline2\nline3")
	}
}

func TestEditorSimulatedTransform(t *testing.T) {
	in, out := simulatedIO("HELLO\n")
	result, err := Editor("Description?",
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
