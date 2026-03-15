package inquire

import (
	"strings"
	"unicode"
)

// TransformTitle converts the first letter of each word to uppercase.
func TransformTitle(s string) string {
	prev := ' '
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(rune(prev)) {
			prev = r
			return unicode.ToTitle(r)
		}
		prev = r
		return r
	}, s)
}

// TransformToLower converts the string to lowercase.
func TransformToLower(s string) string {
	return strings.ToLower(s)
}

// TransformToUpper converts the string to uppercase.
func TransformToUpper(s string) string {
	return strings.ToUpper(s)
}

// ComposeTransformers combines multiple transformers into one.
// Transformers run in order.
func ComposeTransformers(transformers ...func(string) string) func(string) string {
	return func(s string) string {
		for _, t := range transformers {
			s = t(s)
		}
		return s
	}
}
