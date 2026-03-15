package inquire

import (
	"fmt"
	"regexp"
	"strings"
)

// Required returns a validator that rejects empty strings.
func Required(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

// MinLength returns a validator that rejects strings shorter than n.
func MinLength(n int) func(string) error {
	return func(s string) error {
		if len(s) < n {
			return fmt.Errorf("must be at least %d characters", n)
		}
		return nil
	}
}

// MaxLength returns a validator that rejects strings longer than n.
func MaxLength(n int) func(string) error {
	return func(s string) error {
		if len(s) > n {
			return fmt.Errorf("must be at most %d characters", n)
		}
		return nil
	}
}

// MinItems returns a validator for comma-separated values requiring at least n items.
func MinItems(n int) func(string) error {
	return func(s string) error {
		items := countItems(s)
		if items < n {
			return fmt.Errorf("must select at least %d items", n)
		}
		return nil
	}
}

// MaxItems returns a validator for comma-separated values allowing at most n items.
func MaxItems(n int) func(string) error {
	return func(s string) error {
		items := countItems(s)
		if items > n {
			return fmt.Errorf("must select at most %d items", n)
		}
		return nil
	}
}

// MatchRegex returns a validator that requires the input to match a regex pattern.
func MatchRegex(pattern string) func(string) error {
	re := regexp.MustCompile(pattern)
	return func(s string) error {
		if !re.MatchString(s) {
			return fmt.Errorf("must match pattern %s", pattern)
		}
		return nil
	}
}

// ComposeValidators combines multiple validators into one.
// Validators run in order; the first error is returned.
func ComposeValidators(validators ...func(string) error) func(string) error {
	return func(s string) error {
		for _, v := range validators {
			if err := v(s); err != nil {
				return err
			}
		}
		return nil
	}
}

func countItems(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	parts := strings.Split(s, ",")
	count := 0
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			count++
		}
	}
	return count
}
