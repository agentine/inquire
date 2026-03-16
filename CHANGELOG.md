# Changelog

## v0.1.0 — 2026-03-16

Initial release.

### Features

- **Core prompts:** Input, Password, Confirm, Select, MultiSelect, Multiline, Editor
- **Functional options:** WithDefault, WithValidate, WithHelp, WithTransform, WithFilter, WithPageSize, WithContext, WithStdio, WithIcons
- **Structured Ask API:** Fill structs via `Ask([]*Question, &response)` with name-based field mapping
- **Validators:** Required, MinLength, MaxLength, MinItems, MaxItems, MatchRegex, ComposeValidators
- **Transformers:** TransformTitle, TransformToLower, TransformToUpper, ComposeTransformers
- **Terminal I/O:** Raw mode (POSIX termios + Windows Console API), ANSI escape sequences, line editor with cursor movement and word deletion
- **Context support:** All prompts accept `context.Context` for cancellation and timeouts
- **Zero dependencies:** Pure Go, no external packages
- **Compatibility layers:**
  - `compat/survey` — drop-in adapter for AlecAivazis/survey/v2 API (Ask, AskOne, validators, transformers, icons)
  - `compat/promptui` — drop-in adapter for manifoldco/promptui API (Prompt, Select structs with Run)
- **CI:** Test matrix on Go 1.22 + 1.23, linting, publish workflow with pkg.go.dev indexing
