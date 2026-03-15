package inquire

import (
	"errors"
	"strings"

	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// MultiSelect prompts the user to select multiple options from a list.
// Returns the selected indices, the selected strings, and any error.
func MultiSelect(message string, options []string, opts ...Option) ([]int, []string, error) {
	o := applyOptions(opts)

	pio, err := newPromptIO(o)
	if err != nil {
		return nil, nil, err
	}
	defer pio.close()

	if len(options) == 0 {
		return nil, nil, errors.New("inquire: no options provided")
	}

	cursor := 0
	selected := make(map[int]bool)
	pageSize := o.pageSize

	// Pre-select defaults if provided
	if o.defaultVal != nil {
		if defaults, ok := o.defaultVal.([]int); ok {
			for _, idx := range defaults {
				if idx >= 0 && idx < len(options) {
					selected[idx] = true
				}
			}
		}
	}

	filter := ""
	filterFunc := o.filter
	if filterFunc == nil {
		filterFunc = defaultFilter
	}

	filteredIndices := func() []int {
		if filter == "" {
			indices := make([]int, len(options))
			for i := range options {
				indices[i] = i
			}
			return indices
		}
		var indices []int
		for i, opt := range options {
			if filterFunc(filter, opt, i) {
				indices = append(indices, i)
			}
		}
		return indices
	}

	renderMulti := func() {
		indices := filteredIndices()
		term.ClearLine(pio.out)
		pio.write("\r")
		pio.writef("%s%s%s %s ", term.Bold+term.FgCyan, o.icons.Question, term.Reset, message)
		if filter != "" {
			pio.write(filter)
		}
		pio.write("\n")

		start := 0
		if len(indices) > pageSize {
			if cursor >= start+pageSize {
				start = cursor - pageSize + 1
			}
			if cursor < start {
				start = cursor
			}
		}
		end := start + pageSize
		if end > len(indices) {
			end = len(indices)
		}

		for i := start; i < end; i++ {
			term.ClearLine(pio.out)
			idx := indices[i]
			check := "[ ]"
			if selected[idx] {
				check = "[x]"
			}
			if i == cursor {
				pio.writef("  %s%s %s %s%s\n", term.FgCyan, o.icons.Select, check, options[idx], term.Reset)
			} else {
				pio.writef("    %s %s\n", check, options[idx])
			}
		}

		linesWritten := end - start + 1
		term.CursorUp(pio.out, linesWritten)
		term.CursorColumn(pio.out, len(o.icons.Question)+2+len(message)+1+len(filter)+1)
	}

	term.HideCursor(pio.out)
	renderMulti()

	for {
		select {
		case <-o.ctx.Done():
			term.ShowCursor(pio.out)
			pio.write("\n")
			return nil, nil, o.ctx.Err()
		default:
		}

		k, err := pio.reader.ReadKey()
		if err != nil {
			term.ShowCursor(pio.out)
			pio.write("\n")
			return nil, nil, err
		}

		indices := filteredIndices()

		switch k.Type {
		case keys.KeyCtrlC:
			term.ShowCursor(pio.out)
			clearSelectLines(pio, len(indices), pageSize)
			return nil, nil, ErrInterrupt
		case keys.KeyUp, keys.KeyCtrlP:
			if cursor > 0 {
				cursor--
			}
			renderMulti()
		case keys.KeyDown, keys.KeyCtrlN:
			if cursor < len(indices)-1 {
				cursor++
			}
			renderMulti()
		case keys.KeySpace:
			if len(indices) > 0 {
				idx := indices[cursor]
				selected[idx] = !selected[idx]
			}
			renderMulti()
		case keys.KeyEnter:
			var selIndices []int
			var selStrings []string
			for i, opt := range options {
				if selected[i] {
					selIndices = append(selIndices, i)
					selStrings = append(selStrings, opt)
				}
			}

			if o.validate != nil {
				combined := strings.Join(selStrings, ", ")
				if verr := o.validate(combined); verr != nil {
					continue
				}
			}

			term.ShowCursor(pio.out)
			clearSelectLines(pio, len(indices), pageSize)
			answer := strings.Join(selStrings, ", ")
			pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, answer, term.Reset)
			return selIndices, selStrings, nil
		case keys.KeyBackspace, keys.KeyCtrlH:
			if len(filter) > 0 {
				filter = filter[:len(filter)-1]
				cursor = 0
				renderMulti()
			}
		case keys.KeyRune:
			filter += string(k.Rune)
			cursor = 0
			renderMulti()
		}
	}
}
