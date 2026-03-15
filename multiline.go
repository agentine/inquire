package inquire

import (
	"strings"

	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// Multiline prompts for multi-line text input.
// The user presses Enter to add new lines and Ctrl+D or Escape to finish.
func Multiline(message string, opts ...Option) (string, error) {
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

	var lines []string
	if defStr != "" {
		lines = strings.Split(defStr, "\n")
	} else {
		lines = []string{""}
	}
	lineIdx := len(lines) - 1
	colIdx := len(lines[lineIdx])

	renderMultiline := func(errMsg string) {
		term.ClearLine(pio.out)
		pio.write("\r")
		pio.writef("%s%s%s %s %s(Ctrl+D to finish)%s\n", term.Bold+term.FgCyan, o.icons.Question, term.Reset, message, term.Dim, term.Reset)
		for i, line := range lines {
			term.ClearLine(pio.out)
			if i == lineIdx {
				pio.writef("  %s", line)
			} else {
				pio.writef("  %s", line)
			}
			if i < len(lines)-1 {
				pio.write("\n")
			}
		}
		if errMsg != "" {
			pio.writef("\n%s%s %s%s", term.FgRed, o.icons.Error, errMsg, term.Reset)
			term.CursorUp(pio.out, 1)
		}
		// Move cursor to correct position
		linesBelow := len(lines) - 1 - lineIdx
		if linesBelow > 0 {
			term.CursorUp(pio.out, linesBelow)
		}
		term.CursorColumn(pio.out, colIdx+3) // +2 for "  " prefix, +1 for 1-based
	}

	renderMultiline("")

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
		case keys.KeyCtrlD:
			result := strings.Join(lines, "\n")
			result = strings.TrimRight(result, "\n")
			if result == "" && defStr != "" {
				result = defStr
			}
			if o.validate != nil {
				if verr := o.validate(result); verr != nil {
					renderMultiline(verr.Error())
					continue
				}
			}
			if o.transform != nil {
				result = o.transform(result)
			}
			// Clear and show final
			for i := lineIdx; i < len(lines)-1; i++ {
				pio.write("\n")
			}
			pio.write("\n")
			// Move up to header and clear everything
			term.CursorUp(pio.out, len(lines))
			for i := 0; i <= len(lines); i++ {
				term.ClearLine(pio.out)
				if i < len(lines) {
					pio.write("\n")
				}
			}
			term.CursorUp(pio.out, len(lines))
			pio.write("\r")
			preview := result
			if len(preview) > 40 {
				preview = preview[:40] + "..."
			}
			preview = strings.ReplaceAll(preview, "\n", "\\n")
			pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, preview, term.Reset)
			return result, nil
		case keys.KeyEnter:
			// Split current line at cursor
			after := lines[lineIdx][colIdx:]
			lines[lineIdx] = lines[lineIdx][:colIdx]
			lineIdx++
			lines = append(lines[:lineIdx], append([]string{after}, lines[lineIdx:]...)...)
			colIdx = 0
			renderMultiline("")
		case keys.KeyBackspace, keys.KeyCtrlH:
			if colIdx > 0 {
				line := []rune(lines[lineIdx])
				lines[lineIdx] = string(append(line[:colIdx-1], line[colIdx:]...))
				colIdx--
			} else if lineIdx > 0 {
				// Merge with previous line
				prevLen := len(lines[lineIdx-1])
				lines[lineIdx-1] += lines[lineIdx]
				lines = append(lines[:lineIdx], lines[lineIdx+1:]...)
				lineIdx--
				colIdx = prevLen
			}
			renderMultiline("")
		case keys.KeyLeft, keys.KeyCtrlB:
			if colIdx > 0 {
				colIdx--
			} else if lineIdx > 0 {
				lineIdx--
				colIdx = len(lines[lineIdx])
			}
			renderMultiline("")
		case keys.KeyRight, keys.KeyCtrlF:
			if colIdx < len(lines[lineIdx]) {
				colIdx++
			} else if lineIdx < len(lines)-1 {
				lineIdx++
				colIdx = 0
			}
			renderMultiline("")
		case keys.KeyUp, keys.KeyCtrlP:
			if lineIdx > 0 {
				lineIdx--
				if colIdx > len(lines[lineIdx]) {
					colIdx = len(lines[lineIdx])
				}
			}
			renderMultiline("")
		case keys.KeyDown, keys.KeyCtrlN:
			if lineIdx < len(lines)-1 {
				lineIdx++
				if colIdx > len(lines[lineIdx]) {
					colIdx = len(lines[lineIdx])
				}
			}
			renderMultiline("")
		case keys.KeyHome, keys.KeyCtrlA:
			colIdx = 0
			renderMultiline("")
		case keys.KeyEnd, keys.KeyCtrlE:
			colIdx = len(lines[lineIdx])
			renderMultiline("")
		case keys.KeyRune, keys.KeySpace:
			line := []rune(lines[lineIdx])
			newLine := make([]rune, 0, len(line)+1)
			newLine = append(newLine, line[:colIdx]...)
			newLine = append(newLine, k.Rune)
			newLine = append(newLine, line[colIdx:]...)
			lines[lineIdx] = string(newLine)
			colIdx++
			renderMultiline("")
		}
	}
}
