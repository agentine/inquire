package inquire

import (
	"strings"

	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// Confirm prompts for a yes/no answer.
func Confirm(message string, opts ...Option) (bool, error) {
	o := applyOptions(opts)

	pio, err := newPromptIO(o)
	if err != nil {
		return false, err
	}
	defer pio.close()

	defBool := false
	if o.defaultVal != nil {
		if b, ok := o.defaultVal.(bool); ok {
			defBool = b
		}
	}

	hint := "(y/N)"
	if defBool {
		hint = "(Y/n)"
	}

	renderConfirm := func(input string) {
		term.ClearLine(pio.out)
		pio.write("\r")
		pio.writef("%s%s%s %s %s%s%s ", term.Bold+term.FgCyan, o.icons.Question, term.Reset, message, term.Dim, hint, term.Reset)
		if input != "" {
			pio.write(input)
		}
	}

	input := ""
	renderConfirm(input)

	for {
		select {
		case <-o.ctx.Done():
			pio.write("\n")
			return false, o.ctx.Err()
		default:
		}

		k, err := pio.reader.ReadKey()
		if err != nil {
			pio.write("\n")
			return false, err
		}

		switch k.Type {
		case keys.KeyCtrlC:
			pio.write("\n")
			return false, ErrInterrupt
		case keys.KeyEnter:
			result := defBool
			s := strings.TrimSpace(strings.ToLower(input))
			if s == "y" || s == "yes" {
				result = true
			} else if s == "n" || s == "no" {
				result = false
			} else if s != "" {
				// Invalid input — ignore and let them try again
				input = ""
				renderConfirm(input)
				continue
			}

			answer := "No"
			if result {
				answer = "Yes"
			}
			term.ClearLine(pio.out)
			pio.write("\r")
			pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, answer, term.Reset)
			return result, nil
		case keys.KeyBackspace, keys.KeyCtrlH:
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
			renderConfirm(input)
		case keys.KeyRune, keys.KeySpace:
			input += string(k.Rune)
			renderConfirm(input)
		}
	}
}
