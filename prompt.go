package inquire

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// Prompt is the interface that all prompt types implement.
type Prompt interface {
	Run() (any, error)
}

// promptIO resolves the I/O streams to use. If WithStdio was provided, those
// streams are used directly (no raw mode). Otherwise, stdin/stdout are used
// and raw mode is enabled on the terminal.
type promptIO struct {
	in       io.Reader
	out      io.Writer
	reader   *keys.Reader
	rawState *term.State
	isRaw    bool
}

func newPromptIO(o *options) (*promptIO, error) {
	p := &promptIO{}

	if o.in != nil {
		// Simulated I/O — no raw mode
		p.in = o.in
		p.out = o.out
		if p.out == nil {
			p.out = io.Discard
		}
		p.reader = keys.NewReader(p.in)
		return p, nil
	}

	// Real terminal
	p.in = os.Stdin
	p.out = os.Stdout
	if o.out != nil {
		p.out = o.out
	}

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		state, err := term.MakeRaw(fd)
		if err != nil {
			return nil, fmt.Errorf("inquire: failed to set raw mode: %w", err)
		}
		p.rawState = state
		p.isRaw = true
	}

	p.reader = keys.NewReader(p.in)
	return p, nil
}

func (p *promptIO) close() {
	if p.rawState != nil {
		_ = term.Restore(p.rawState)
	}
}

// write writes a string to the output, translating \n to \r\n in raw mode.
func (p *promptIO) write(s string) {
	if p.isRaw {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	fmt.Fprint(p.out, s)
}

// writef writes a formatted string.
func (p *promptIO) writef(format string, args ...any) {
	p.write(fmt.Sprintf(format, args...))
}
