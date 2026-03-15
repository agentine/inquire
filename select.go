package inquire

import (
	"errors"
	"strings"

	"github.com/agentine/inquire/internal/keys"
	"github.com/agentine/inquire/internal/term"
)

// Select prompts the user to select one option from a list.
// Returns the selected index, the selected string, and any error.
func Select(message string, options []string, opts ...Option) (int, string, error) {
	o := applyOptions(opts)

	pio, err := newPromptIO(o)
	if err != nil {
		return -1, "", err
	}
	defer pio.close()

	if len(options) == 0 {
		return -1, "", errors.New("inquire: no options provided")
	}

	cursor := 0
	if o.defaultVal != nil {
		switch v := o.defaultVal.(type) {
		case int:
			if v >= 0 && v < len(options) {
				cursor = v
			}
		case string:
			for i, opt := range options {
				if opt == v {
					cursor = i
					break
				}
			}
		}
	}

	filter := ""
	pageSize := o.pageSize

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

	renderSelect := func() {
		indices := filteredIndices()
		term.ClearLine(pio.out)
		pio.write("\r")
		pio.writef("%s%s%s %s ", term.Bold+term.FgCyan, o.icons.Question, term.Reset, message)
		if filter != "" {
			pio.write(filter)
		}
		pio.write("\n")

		// Pagination
		start := 0
		if len(indices) > pageSize {
			// Keep cursor visible
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
			if i == cursor {
				pio.writef("  %s%s %s%s\n", term.FgCyan, o.icons.Select, options[idx], term.Reset)
			} else {
				pio.writef("    %s\n", options[idx])
			}
		}

		// Move cursor back to prompt line
		linesWritten := end - start + 1
		term.CursorUp(pio.out, linesWritten)
		// Position cursor at end of filter text
		term.CursorColumn(pio.out, len(o.icons.Question)+2+len(message)+1+len(filter)+1)
	}

	term.HideCursor(pio.out)
	renderSelect()

	for {
		select {
		case <-o.ctx.Done():
			term.ShowCursor(pio.out)
			pio.write("\n")
			return -1, "", o.ctx.Err()
		default:
		}

		k, err := pio.reader.ReadKey()
		if err != nil {
			term.ShowCursor(pio.out)
			pio.write("\n")
			return -1, "", err
		}

		indices := filteredIndices()

		switch k.Type {
		case keys.KeyCtrlC:
			term.ShowCursor(pio.out)
			clearSelectLines(pio, len(indices), pageSize)
			return -1, "", ErrInterrupt
		case keys.KeyUp, keys.KeyCtrlP:
			if cursor > 0 {
				cursor--
			}
			renderSelect()
		case keys.KeyDown, keys.KeyCtrlN:
			if cursor < len(indices)-1 {
				cursor++
			}
			renderSelect()
		case keys.KeyEnter:
			if len(indices) == 0 {
				continue
			}
			selectedIdx := indices[cursor]
			selected := options[selectedIdx]
			term.ShowCursor(pio.out)
			clearSelectLines(pio, len(indices), pageSize)
			pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, selected, term.Reset)
			return selectedIdx, selected, nil
		case keys.KeyBackspace, keys.KeyCtrlH:
			if len(filter) > 0 {
				filter = filter[:len(filter)-1]
				cursor = 0
				renderSelect()
			}
		case keys.KeyRune, keys.KeySpace:
			filter += string(k.Rune)
			cursor = 0
			renderSelect()
		}
	}
}

func clearSelectLines(pio *promptIO, numItems, pageSize int) {
	visible := numItems
	if visible > pageSize {
		visible = pageSize
	}
	// Clear the option lines below
	for i := 0; i <= visible; i++ {
		term.ClearLine(pio.out)
		if i < visible {
			pio.write("\n")
		}
	}
	term.CursorUp(pio.out, visible)
	pio.write("\r")
}

func defaultFilter(filter, option string, _ int) bool {
	return strings.Contains(strings.ToLower(option), strings.ToLower(filter))
}
