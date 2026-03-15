package inquire

import (
	"errors"
	"testing"
)

func TestAskBasic(t *testing.T) {
	in1, out := simulatedIO("Alice\n")
	in2, _ := simulatedIO("secret\n")
	in3, _ := simulatedIO("y\n")
	var result struct {
		Name     string
		Password string
		Confirm  bool
	}

	err := Ask([]*Question{
		{Name: "Name", Prompt: &InputPrompt{Message: "Name?", Options: []Option{WithStdio(in1, out, out)}}},
		{Name: "Password", Prompt: &PasswordPrompt{Message: "Password?", Options: []Option{WithStdio(in2, out, out)}}},
		{Name: "Confirm", Prompt: &ConfirmPrompt{Message: "OK?", Options: []Option{WithStdio(in3, out, out)}}},
	}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", result.Name)
	}
	if result.Password != "secret" {
		t.Errorf("Password = %q, want secret", result.Password)
	}
	if !result.Confirm {
		t.Error("Confirm = false, want true")
	}
}

func TestAskNotPointer(t *testing.T) {
	var result struct{ Name string }
	err := Ask(nil, result) // not a pointer
	if err == nil {
		t.Error("expected error for non-pointer")
	}
}

func TestAskInvalidField(t *testing.T) {
	in, out := simulatedIO("test\n")
	var result struct{ Name string }
	err := Ask([]*Question{
		{Name: "Missing", Prompt: &InputPrompt{Message: "?", Options: []Option{WithStdio(in, out, out)}}},
	}, &result)
	if err == nil {
		t.Error("expected error for missing field")
	}
}

func TestAskWithValidation(t *testing.T) {
	in, out := simulatedIO("test\n")
	var result struct{ Name string }
	err := Ask([]*Question{
		{
			Name:   "Name",
			Prompt: &InputPrompt{Message: "Name?", Options: []Option{WithStdio(in, out, out)}},
			Validate: func(v any) error {
				s, _ := v.(string)
				if s == "" {
					return errors.New("required")
				}
				return nil
			},
		},
	}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "test" {
		t.Errorf("Name = %q, want test", result.Name)
	}
}

func TestAskWithTransform(t *testing.T) {
	in, out := simulatedIO("HELLO\n")
	var result struct{ Name string }
	err := Ask([]*Question{
		{
			Name:   "Name",
			Prompt: &InputPrompt{Message: "Name?", Options: []Option{WithStdio(in, out, out)}},
			Transform: func(v any) any {
				s, _ := v.(string)
				return TransformToLower(s)
			},
		},
	}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Name != "hello" {
		t.Errorf("Name = %q, want hello", result.Name)
	}
}

func TestAskSelect(t *testing.T) {
	// Press enter to select first item
	in, out := simulatedIO("\n")
	var result struct{ Color string }
	err := Ask([]*Question{
		{
			Name:   "Color",
			Prompt: &SelectPrompt{Message: "Color?", Items: []string{"red", "green"}, Options: []Option{WithStdio(in, out, out)}},
		},
	}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Color != "red" {
		t.Errorf("Color = %q, want red", result.Color)
	}
}
