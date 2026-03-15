package inquire

import (
	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// Input prompts for a single line of text input.
func Input(message string, opts ...Option) (string, error) {
	o := applyOptions(opts)

	pio, err := newPromptIO(o)
	if err != nil {
		return "", err
	}
	defer pio.close()

	defStr := ""
	if o.defaultVal != nil {
		if s, ok := o.defaultVal.(string); ok {
			defStr = s
		}
	}

	line := keys.NewLine()

	renderInput := func(errMsg string) {
		term.ClearLine(pio.out)
		pio.write("\r")
		pio.writef("%s%s%s %s ", term.Bold+term.FgCyan, o.icons.Question, term.Reset, message)
		if line.Len() == 0 && defStr != "" {
			pio.writef("%s(%s)%s ", term.Dim, defStr, term.Reset)
		}
		pio.write(line.String())
		// Move cursor to correct position
		if line.Pos() < line.Len() {
			term.CursorBack(pio.out, line.Len()-line.Pos())
		}
		if errMsg != "" {
			// Save position, write error below, restore
			term.SaveCursor(pio.out)
			pio.writef("\n%s%s %s%s", term.FgRed, o.icons.Error, errMsg, term.Reset)
			term.RestoreCursor(pio.out)
		}
	}

	renderInput("")

	for {
		// Check context
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
			if result == "" && defStr != "" {
				result = defStr
			}
			if o.validate != nil {
				if verr := o.validate(result); verr != nil {
					renderInput(verr.Error())
					continue
				}
			}
			if o.transform != nil {
				result = o.transform(result)
			}
			// Clear error line if any and show final answer
			term.ClearLine(pio.out)
			pio.write("\r")
			pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, result, term.Reset)
			return result, nil
		case keys.KeyBackspace, keys.KeyCtrlH:
			line.Backspace()
			renderInput("")
		case keys.KeyDelete:
			line.Delete()
			renderInput("")
		case keys.KeyLeft, keys.KeyCtrlB:
			line.MoveLeft()
			renderInput("")
		case keys.KeyRight, keys.KeyCtrlF:
			line.MoveRight()
			renderInput("")
		case keys.KeyHome, keys.KeyCtrlA:
			line.MoveHome()
			renderInput("")
		case keys.KeyEnd, keys.KeyCtrlE:
			line.MoveEnd()
			renderInput("")
		case keys.KeyCtrlK:
			line.DeleteToEnd()
			renderInput("")
		case keys.KeyCtrlU:
			line.DeleteToStart()
			renderInput("")
		case keys.KeyCtrlW:
			line.DeleteWord()
			renderInput("")
		case keys.KeyRune, keys.KeySpace:
			line.Insert(k.Rune)
			renderInput("")
		}
	}
}
