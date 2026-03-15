package inquire

import "testing"

func TestTransformTitle(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"hello world", "Hello World"},
		{"HELLO", "HELLO"},
		{"", ""},
		{"a", "A"},
	}
	for _, tt := range tests {
		if got := TransformTitle(tt.in); got != tt.want {
			t.Errorf("TransformTitle(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestTransformToLower(t *testing.T) {
	if got := TransformToLower("HELLO"); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestTransformToUpper(t *testing.T) {
	if got := TransformToUpper("hello"); got != "HELLO" {
		t.Errorf("got %q, want %q", got, "HELLO")
	}
}

func TestComposeTransformers(t *testing.T) {
	tf := ComposeTransformers(TransformToLower, TransformTitle)
	if got := tf("HELLO WORLD"); got != "Hello World" {
		t.Errorf("got %q, want %q", got, "Hello World")
	}
}
