package inquire

import "testing"

func TestRequired(t *testing.T) {
	if err := Required("hello"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := Required(""); err == nil {
		t.Error("expected error for empty string")
	}
	if err := Required("   "); err == nil {
		t.Error("expected error for whitespace-only string")
	}
}

func TestMinLength(t *testing.T) {
	v := MinLength(3)
	if err := v("ab"); err == nil {
		t.Error("expected error for short string")
	}
	if err := v("abc"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v("abcd"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMaxLength(t *testing.T) {
	v := MaxLength(5)
	if err := v("abc"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v("abcde"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v("abcdef"); err == nil {
		t.Error("expected error for long string")
	}
}

func TestMinItems(t *testing.T) {
	v := MinItems(2)
	if err := v("a, b"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v("a"); err == nil {
		t.Error("expected error for too few items")
	}
	if err := v(""); err == nil {
		t.Error("expected error for empty")
	}
}

func TestMaxItems(t *testing.T) {
	v := MaxItems(2)
	if err := v("a, b"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v("a, b, c"); err == nil {
		t.Error("expected error for too many items")
	}
}

func TestMatchRegex(t *testing.T) {
	v := MatchRegex(`^\d{3}-\d{4}$`)
	if err := v("123-4567"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v("abc"); err == nil {
		t.Error("expected error for non-matching string")
	}
}

func TestComposeValidators(t *testing.T) {
	v := ComposeValidators(Required, MinLength(3))
	if err := v(""); err == nil {
		t.Error("expected error for empty")
	}
	if err := v("ab"); err == nil {
		t.Error("expected error for short")
	}
	if err := v("abc"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
