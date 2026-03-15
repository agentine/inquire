// Package survey provides drop-in compatibility with AlecAivazis/survey/v2.
//
// Migration guide:
//
// 1. Replace import "github.com/AlecAivazis/survey/v2" with
//    "github.com/agentine/inquire/compat/survey"
// 2. Recompile. No other changes needed.
//
// This package maps survey/v2 types and functions to their inquire equivalents.
package survey

import (
	"io"

	"github.com/agentine/inquire"
)

// Prompt is the interface that all prompt types implement.
type Prompt interface {
	prompt() inquire.Prompt
}

// AskOpt configures the Ask/AskOne call.
type AskOpt func(*askOptions)

type askOptions struct {
	stdio *stdioOpt
	icons *inquire.IconSet
}

type stdioOpt struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

// WithStdio overrides the I/O streams.
func WithStdio(in io.Reader, out io.Writer, err io.Writer) AskOpt {
	return func(o *askOptions) {
		o.stdio = &stdioOpt{in: in, out: out, err: err}
	}
}

// WithIcons overrides the icon set.
func WithIcons(fn func(*IconSet)) AskOpt {
	return func(o *askOptions) {
		is := DefaultIcons()
		fn(is)
		o.icons = &inquire.IconSet{
			Question: is.Question.Text,
			Help:     is.Help.Text,
			Error:    is.Error.Text,
			Select:   is.SelectFocus.Text,
		}
	}
}

func resolveOpts(askOpts []AskOpt) []inquire.Option {
	ao := &askOptions{}
	for _, fn := range askOpts {
		fn(ao)
	}
	var opts []inquire.Option
	if ao.stdio != nil {
		opts = append(opts, inquire.WithStdio(ao.stdio.in, ao.stdio.out, ao.stdio.err))
	}
	if ao.icons != nil {
		opts = append(opts, inquire.WithIcons(*ao.icons))
	}
	return opts
}

// Input is a single-line text input prompt.
type Input struct {
	Message string
	Default string
	Help    string
}

func (i *Input) prompt() inquire.Prompt {
	var opts []inquire.Option
	if i.Default != "" {
		opts = append(opts, inquire.WithDefault(i.Default))
	}
	if i.Help != "" {
		opts = append(opts, inquire.WithHelp(i.Help))
	}
	return &inquire.InputPrompt{Message: i.Message, Options: opts}
}

// Password is a masked input prompt.
type Password struct {
	Message string
	Help    string
}

func (p *Password) prompt() inquire.Prompt {
	var opts []inquire.Option
	if p.Help != "" {
		opts = append(opts, inquire.WithHelp(p.Help))
	}
	return &inquire.PasswordPrompt{Message: p.Message, Options: opts}
}

// Confirm is a yes/no prompt.
type Confirm struct {
	Message string
	Default bool
	Help    string
}

func (c *Confirm) prompt() inquire.Prompt {
	var opts []inquire.Option
	opts = append(opts, inquire.WithDefault(c.Default))
	if c.Help != "" {
		opts = append(opts, inquire.WithHelp(c.Help))
	}
	return &inquire.ConfirmPrompt{Message: c.Message, Options: opts}
}

// Select is a single-selection prompt.
type Select struct {
	Message  string
	Options  []string
	Default  interface{}
	Help     string
	PageSize int
}

func (s *Select) prompt() inquire.Prompt {
	var opts []inquire.Option
	if s.Default != nil {
		opts = append(opts, inquire.WithDefault(s.Default))
	}
	if s.Help != "" {
		opts = append(opts, inquire.WithHelp(s.Help))
	}
	if s.PageSize > 0 {
		opts = append(opts, inquire.WithPageSize(s.PageSize))
	}
	return &inquire.SelectPrompt{Message: s.Message, Items: s.Options, Options: opts}
}

// MultiSelect is a multi-selection prompt.
type MultiSelect struct {
	Message  string
	Options  []string
	Default  []string
	Help     string
	PageSize int
}

func (m *MultiSelect) prompt() inquire.Prompt {
	var opts []inquire.Option
	if m.Default != nil {
		// Convert default strings to indices
		var indices []int
		for _, def := range m.Default {
			for i, opt := range m.Options {
				if opt == def {
					indices = append(indices, i)
					break
				}
			}
		}
		if len(indices) > 0 {
			opts = append(opts, inquire.WithDefault(indices))
		}
	}
	if m.Help != "" {
		opts = append(opts, inquire.WithHelp(m.Help))
	}
	if m.PageSize > 0 {
		opts = append(opts, inquire.WithPageSize(m.PageSize))
	}
	return &inquire.MultiSelectPrompt{Message: m.Message, Items: m.Options, Options: opts}
}

// Multiline is a multi-line text input prompt.
type Multiline struct {
	Message string
	Default string
	Help    string
}

func (m *Multiline) prompt() inquire.Prompt {
	var opts []inquire.Option
	if m.Default != "" {
		opts = append(opts, inquire.WithDefault(m.Default))
	}
	if m.Help != "" {
		opts = append(opts, inquire.WithHelp(m.Help))
	}
	return &inquire.MultilinePrompt{Message: m.Message, Options: opts}
}

// Editor opens an external editor.
type Editor struct {
	Message string
	Default string
	Help    string
}

func (e *Editor) prompt() inquire.Prompt {
	var opts []inquire.Option
	if e.Default != "" {
		opts = append(opts, inquire.WithDefault(e.Default))
	}
	if e.Help != "" {
		opts = append(opts, inquire.WithHelp(e.Help))
	}
	return &inquire.EditorPrompt{Message: e.Message, Options: opts}
}

// Validator is a function that validates input.
type Validator func(ans interface{}) error

// Transformer is a function that transforms input.
type Transformer func(ans interface{}) interface{}

// Required validates that the input is not empty.
func Required(ans interface{}) error {
	s, _ := ans.(string)
	return inquire.Required(s)
}

// MinLength returns a validator that rejects strings shorter than n.
func MinLength(n int) Validator {
	v := inquire.MinLength(n)
	return func(ans interface{}) error {
		s, _ := ans.(string)
		return v(s)
	}
}

// MaxLength returns a validator that rejects strings longer than n.
func MaxLength(n int) Validator {
	v := inquire.MaxLength(n)
	return func(ans interface{}) error {
		s, _ := ans.(string)
		return v(s)
	}
}

// MinItems returns a validator requiring at least n selected items.
func MinItems(n int) Validator {
	v := inquire.MinItems(n)
	return func(ans interface{}) error {
		s, _ := ans.(string)
		return v(s)
	}
}

// MaxItems returns a validator allowing at most n selected items.
func MaxItems(n int) Validator {
	v := inquire.MaxItems(n)
	return func(ans interface{}) error {
		s, _ := ans.(string)
		return v(s)
	}
}

// ComposeValidators combines multiple validators.
func ComposeValidators(validators ...Validator) Validator {
	return func(ans interface{}) error {
		for _, v := range validators {
			if err := v(ans); err != nil {
				return err
			}
		}
		return nil
	}
}

// Title transforms the input to title case.
func Title(ans interface{}) interface{} {
	s, _ := ans.(string)
	return inquire.TransformTitle(s)
}

// ToLower transforms the input to lowercase.
func ToLower(ans interface{}) interface{} {
	s, _ := ans.(string)
	return inquire.TransformToLower(s)
}

// ToUpper transforms the input to uppercase.
func ToUpper(ans interface{}) interface{} {
	s, _ := ans.(string)
	return inquire.TransformToUpper(s)
}

// ComposeTransformers combines multiple transformers.
func ComposeTransformers(transformers ...Transformer) Transformer {
	return func(ans interface{}) interface{} {
		for _, t := range transformers {
			ans = t(ans)
		}
		return ans
	}
}

// Question defines a single question for Ask.
type Question struct {
	Name      string
	Prompt    Prompt
	Validate  Validator
	Transform Transformer
}

// Ask runs a set of questions and fills the response struct.
func Ask(qs []*Question, response interface{}, opts ...AskOpt) error {
	baseOpts := resolveOpts(opts)

	questions := make([]*inquire.Question, len(qs))
	for i, q := range qs {
		p := q.Prompt.prompt()

		// Inject base options (stdio, icons) into the prompt
		injectOpts(p, baseOpts)

		iq := &inquire.Question{
			Name:   q.Name,
			Prompt: p,
		}
		if q.Validate != nil {
			v := q.Validate
			iq.Validate = func(ans any) error { return v(ans) }
		}
		if q.Transform != nil {
			t := q.Transform
			iq.Transform = func(ans any) any { return t(ans) }
		}
		questions[i] = iq
	}

	return inquire.Ask(questions, response)
}

// AskOne runs a single prompt and stores the result.
func AskOne(p Prompt, response interface{}, opts ...AskOpt) error {
	baseOpts := resolveOpts(opts)
	prompt := p.prompt()
	injectOpts(prompt, baseOpts)

	answer, err := prompt.Run()
	if err != nil {
		return err
	}

	// Set the response value
	return setResponse(response, answer)
}

func injectOpts(p inquire.Prompt, opts []inquire.Option) {
	switch pt := p.(type) {
	case *inquire.InputPrompt:
		pt.Options = append(pt.Options, opts...)
	case *inquire.PasswordPrompt:
		pt.Options = append(pt.Options, opts...)
	case *inquire.ConfirmPrompt:
		pt.Options = append(pt.Options, opts...)
	case *inquire.SelectPrompt:
		pt.Options = append(pt.Options, opts...)
	case *inquire.MultiSelectPrompt:
		pt.Options = append(pt.Options, opts...)
	case *inquire.MultilinePrompt:
		pt.Options = append(pt.Options, opts...)
	case *inquire.EditorPrompt:
		pt.Options = append(pt.Options, opts...)
	}
}
