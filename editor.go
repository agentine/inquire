package inquire

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/agentine/inquire/internal/term"
)

// Editor opens the user's $EDITOR with a temp file and returns the content.
func Editor(message string, opts ...Option) (string, error) {
	o := applyOptions(opts)

	pio, err := newPromptIO(o)
	if err != nil {
		return "", err
	}

	// For simulated I/O (testing), just read all input as the content
	if o.in != nil {
		defer pio.close()
		buf, err := readAll(o.in)
		if err != nil {
			return "", err
		}
		result := strings.TrimRight(string(buf), "\n")
		if o.validate != nil {
			if verr := o.validate(result); verr != nil {
				return "", verr
			}
		}
		if o.transform != nil {
			result = o.transform(result)
		}
		pio.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, truncate(result), term.Reset)
		return result, nil
	}

	// Restore terminal before opening editor
	pio.close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	defStr := ""
	if o.defaultVal != nil {
		if s, ok := o.defaultVal.(string); ok {
			defStr = s
		}
	}

	tmpfile, err := os.CreateTemp("", "inquire-*.txt")
	if err != nil {
		return "", err
	}
	tmpPath := tmpfile.Name()
	defer os.Remove(tmpPath)

	if defStr != "" {
		if _, err := tmpfile.WriteString(defStr); err != nil {
			tmpfile.Close()
			return "", err
		}
	}
	tmpfile.Close()

	for {
		cmd := exec.Command(editor, tmpPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", err
		}

		content, err := os.ReadFile(tmpPath)
		if err != nil {
			return "", err
		}

		result := strings.TrimRight(string(content), "\n")
		if o.validate != nil {
			if verr := o.validate(result); verr != nil {
				// Show error and reopen editor
				var w io.Writer = os.Stdout
				if o.out != nil {
					w = o.out
				}
				pio2 := &promptIO{out: w}
				pio2.writef("%s%s %s%s\n", term.FgRed, o.icons.Error, verr.Error(), term.Reset)
				continue
			}
		}
		if o.transform != nil {
			result = o.transform(result)
		}

		var w io.Writer = os.Stdout
		if o.out != nil {
			w = o.out
		}
		pio2 := &promptIO{out: w}
		pio2.writef("%s%s%s %s %s%s%s\n", term.Bold+term.FgGreen, o.icons.Question, term.Reset, message, term.FgCyan, truncate(result), term.Reset)
		return result, nil
	}
}

func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var buf []byte
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			break
		}
	}
	return buf, nil
}

func truncate(s string) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	if len(s) > 40 {
		return s[:40] + "..."
	}
	return s
}
