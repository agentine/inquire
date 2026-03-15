package inquire

import (
	"context"
	"io"
)

// Option configures a prompt.
type Option func(*options)

// IconSet configures the icons used by prompts.
type IconSet struct {
	Question string
	Help     string
	Error    string
	Select   string
}

type options struct {
	defaultVal any
	help       string
	validate   func(string) error
	transform  func(string) string
	filter     func(filter string, option string, index int) bool
	pageSize   int
	ctx        context.Context
	in         io.Reader
	out        io.Writer
	errOut     io.Writer
	icons      *IconSet
}

func defaultOptions() *options {
	return &options{
		pageSize: 7,
		ctx:      context.Background(),
		icons: &IconSet{
			Question: "?",
			Help:     "i",
			Error:    "X",
			Select:   ">",
		},
	}
}

func applyOptions(opts []Option) *options {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}
	return o
}

// WithDefault sets the default value for a prompt.
func WithDefault(val any) Option {
	return func(o *options) {
		o.defaultVal = val
	}
}

// WithHelp sets help text shown when the user presses '?'.
func WithHelp(msg string) Option {
	return func(o *options) {
		o.help = msg
	}
}

// WithValidate sets a validation function for the input.
func WithValidate(fn func(string) error) Option {
	return func(o *options) {
		o.validate = fn
	}
}

// WithTransform sets a transformation function applied to the final answer.
func WithTransform(fn func(string) string) Option {
	return func(o *options) {
		o.transform = fn
	}
}

// WithFilter sets a custom filter function for Select/MultiSelect.
func WithFilter(fn func(filter string, option string, index int) bool) Option {
	return func(o *options) {
		o.filter = fn
	}
}

// WithPageSize sets the number of items displayed at once in Select/MultiSelect.
func WithPageSize(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.pageSize = n
		}
	}
}

// WithContext sets a context for cancellation/timeout.
func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// WithStdio overrides the input/output streams for the prompt.
func WithStdio(in io.Reader, out io.Writer, errOut io.Writer) Option {
	return func(o *options) {
		o.in = in
		o.out = out
		o.errOut = errOut
	}
}

// WithIcons overrides the icon set used by prompts.
func WithIcons(icons IconSet) Option {
	return func(o *options) {
		o.icons = &icons
	}
}
