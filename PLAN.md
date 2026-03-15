# inquire — Interactive Terminal Prompts for Go

## Overview

Drop-in replacement for **AlecAivazis/survey** (3,089 importers, 4.1k stars, archived April 2024) and **manifoldco/promptui** (3,470 importers, 6.1k stars, unmaintained since Oct 2021). No maintained drop-in fork exists for either library. The recommended alternative (charmbracelet/huh) uses a fundamentally different Elm-style architecture that requires rewriting consumer code.

**Package name:** `github.com/agentine/inquire` (verified available on pkg.go.dev)

## Target Libraries

| Library | Importers | Stars | Last Release | Status |
|---------|-----------|-------|-------------|--------|
| AlecAivazis/survey/v2 | 3,089 | 4.1k | v2.3.7 (Dec 2022) | Archived Apr 2024 |
| manifoldco/promptui | 3,470 | 6.1k | v0.9.0 (Oct 2021) | Unmaintained |

## Architecture

### Core Package (`inquire`)

Primary API — simple, functional prompt calls:

- `Input(message string, opts ...Option) (string, error)` — single-line text input
- `Multiline(message string, opts ...Option) (string, error)` — multi-line text input
- `Password(message string, opts ...Option) (string, error)` — masked input
- `Confirm(message string, opts ...Option) (bool, error)` — yes/no confirmation
- `Select(message string, options []string, opts ...Option) (int, string, error)` — single selection
- `MultiSelect(message string, options []string, opts ...Option) ([]int, []string, error)` — multiple selection
- `Editor(message string, opts ...Option) (string, error)` — external editor

### Functional Options

```go
WithDefault(val any)
WithHelp(msg string)
WithValidate(fn func(string) error)
WithTransform(fn func(string) string)
WithFilter(fn func(filter string, option string, index int) bool)
WithPageSize(n int)
WithContext(ctx context.Context)
WithStdio(in io.Reader, out io.Writer, err io.Writer)
WithIcons(icons IconSet)
```

### Structured Form API

```go
type Question struct {
    Name      string
    Prompt    Prompt
    Validate  func(any) error
    Transform func(any) any
}

func Ask(questions []*Question, response any, opts ...Option) error
```

### Compatibility Layers

- `inquire/compat/survey` — type aliases and adapters for survey/v2 API
- `inquire/compat/promptui` — adapters for promptui API

## Key Improvements Over survey/promptui

1. **Context support** — all prompts accept `context.Context` for cancellation/timeout
2. **Generics** — type-safe select/multiselect with `SelectTyped[T]` / `MultiSelectTyped[T]`
3. **Better Windows support** — proper Windows Console API, no stale terminal state on interrupt
4. **Zero dependencies** — pure Go, no external packages
5. **Accessible** — screen reader support, high-contrast mode
6. **Modern Go** — Go 1.22+ with iterators where appropriate

## Major Components

1. **Terminal I/O** — raw mode, ANSI sequences, cursor control, Windows Console API
2. **Prompt rendering** — template-based rendering with ANSI styling
3. **Input handling** — key events, line editing, history, filtering
4. **Validators** — Required, MinLength, MaxLength, Regex, custom
5. **Transformers** — Title, ToLower, ToUpper, ComposeTransformers
6. **Compatibility** — survey/v2 and promptui adapter packages

## Deliverables

- `inquire` core package with all prompt types
- `inquire/compat/survey` compatibility layer
- `inquire/compat/promptui` compatibility layer
- Comprehensive test suite
- Examples and migration guide
- Go module published at `github.com/agentine/inquire`
