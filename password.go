package inquire

import (
	"strings"

	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// Password prompts for a password with masked input (asterisks).
func Password(message string, opts ...Option) (string, error) {
	o := applyOptions(opts)

	pio, err := newPromptIO(o)
	if err != nil {
		return "", err
	}
	defer pio.close()

	line := keys.NewLine()

	renderPassword := func(errMsg string) {
		term.ClearLine(pio.out)
		pio.write("\r")
		pio.writef("%s%s%s %s ", term.Bold+term.FgCyan, o.icons.Question, term.Reset, message)
		pio.write(strings.Repeat("*", line.Len()))
		if line.Pos() < line.Len() {
			term.CursorBack(pio.out, line.Len()-line.Pos())
		}
		if errMsg != "" {
			term.SaveCursor(pio.out)
			pio.writef("\n%s%s %s%s", term.FgRed, o.icons.Error, errMsg, term.Reset)
			term.RestoreCursor(pio.out)
		}
	}

	renderPassword("")

	for {
		select {
		case <-o.ctx.Done():
			pio.write("\n")
			return "", o.ctx.Err()
		default:
		}

		k, err := pio.reader.ReadKey()
		if err != nil {
			pio.write("\n")
			return "", err
		}

		switch k.Type {
		case keys.KeyCtrlC:
			pio.write("\n")
			return "", ErrInterrupt
		case keys.KeyEnter:
			result := line.String()
			if o.validate != nil {
				if verr := o.validate(result); verr != nil {
					renderPassword(verr.Error())
					continue
				}
			}
			if o.transform != nil {
				result = o.transform(result)
			}
			term.ClearLine(pio.out)
			pio.write("\r")
			pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, strings.Repeat("*", len([]rune(result))), term.Reset)
			return result, nil
		case keys.KeyBackspace, keys.KeyCtrlH:
			line.Backspace()
			renderPassword("")
		case keys.KeyDelete:
			line.Delete()
			renderPassword("")
		case keys.KeyLeft, keys.KeyCtrlB:
			line.MoveLeft()
			renderPassword("")
		case keys.KeyRight, keys.KeyCtrlF:
			line.MoveRight()
			renderPassword("")
		case keys.KeyHome, keys.KeyCtrlA:
			line.MoveHome()
			renderPassword("")
		case keys.KeyEnd, keys.KeyCtrlE:
			line.MoveEnd()
			renderPassword("")
		case keys.KeyCtrlU:
			line.DeleteToStart()
			renderPassword("")
		case keys.KeyCtrlW:
			line.DeleteWord()
			renderPassword("")
		case keys.KeyRune, keys.KeySpace:
			line.Insert(k.Rune)
			renderPassword("")
		}
	}
}
