// Package promptui provides drop-in compatibility with manifoldco/promptui.
//
// Migration guide:
//
// 1. Replace import "github.com/manifoldco/promptui" with
//    "github.com/agentine/inquire/compat/promptui"
// 2. Recompile. No other changes needed.
//
// This package maps promptui types and functions to their inquire equivalents.
package promptui

import (
	"io"
	"strings"

	"github.com/agentine/inquire"
)

// ErrInterrupt is returned when the user presses Ctrl+C.
var ErrInterrupt = inquire.ErrInterrupt

// ErrEOF is returned when the input ends unexpectedly.
var ErrEOF = io.EOF

// ValidateFunc is the type for validation functions.
type ValidateFunc func(string) error

// Prompt represents a single-line text prompt (matching promptui.Prompt).
type Prompt struct {
	Label     interface{}
	Default   string
	Validate  ValidateFunc
	Mask      rune
	IsConfirm bool
	IsVimMode bool
	Pointer   func([]rune, int) []rune // cursor pointer (ignored in compat)
	Stdin     io.ReadCloser
	Stdout    io.WriteCloser
}

// Run executes the prompt and returns the result.
func (p *Prompt) Run() (string, error) {
	label := labelString(p.Label)

	var opts []inquire.Option
	if p.Default != "" {
		opts = append(opts, inquire.WithDefault(p.Default))
	}
	if p.Validate != nil {
		v := p.Validate
		opts = append(opts, inquire.WithValidate(v))
	}
	if p.Stdin != nil || p.Stdout != nil {
		in := io.Reader(p.Stdin)
		out := io.Writer(p.Stdout)
		if in == nil {
			in = strings.NewReader("\n")
		}
		if out == nil {
			out = io.Discard
		}
		opts = append(opts, inquire.WithStdio(in, out, out))
	}

	if p.IsConfirm {
		result, err := inquire.Confirm(label, opts...)
		if err != nil {
			return "", err
		}
		if result {
			return "y", nil
		}
		return "N", nil
	}

	if p.Mask != 0 {
		return inquire.Password(label, opts...)
	}

	return inquire.Input(label, opts...)
}

// Select represents a selection prompt (matching promptui.Select).
type Select struct {
	Label             interface{}
	Items             interface{}
	Size              int
	Templates         *SelectTemplates
	Searcher          func(input string, index int) bool
	StartInSearchMode bool
	Stdin             io.ReadCloser
	Stdout            io.WriteCloser
}

// SelectTemplates holds templates for the select prompt (compatibility stub).
type SelectTemplates struct {
	Label    string
	Active   string
	Inactive string
	Selected string
	Details  string
	FuncMap  interface{}
}

// Run executes the select prompt and returns the index and value.
func (s *Select) Run() (int, string, error) {
	label := labelString(s.Label)
	items := toStringSlice(s.Items)

	var opts []inquire.Option
	if s.Size > 0 {
		opts = append(opts, inquire.WithPageSize(s.Size))
	}
	if s.Searcher != nil {
		searcher := s.Searcher
		opts = append(opts, inquire.WithFilter(func(filter string, option string, index int) bool {
			return searcher(filter, index)
		}))
	}
	if s.Stdin != nil || s.Stdout != nil {
		in := io.Reader(s.Stdin)
		out := io.Writer(s.Stdout)
		if in == nil {
			in = strings.NewReader("\n")
		}
		if out == nil {
			out = io.Discard
		}
		opts = append(opts, inquire.WithStdio(in, out, out))
	}

	return inquire.Select(label, items, opts...)
}

func labelString(label interface{}) string {
	if s, ok := label.(string); ok {
		return s
	}
	if s, ok := label.(interface{ String() string }); ok {
		return s.String()
	}
	return ""
}

func toStringSlice(items interface{}) []string {
	switch v := items.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = labelString(item)
		}
		return result
	default:
		return nil
	}
}
