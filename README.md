# inquire

Interactive terminal prompts for Go.

Drop-in replacement for [AlecAivazis/survey](https://github.com/AlecAivazis/survey) and [manifoldco/promptui](https://github.com/manifoldco/promptui).

## Features

- **7 prompt types** — Input, Password, Confirm, Select, MultiSelect, Multiline, Editor
- **Context support** — all prompts accept `context.Context` for cancellation/timeout
- **Zero dependencies** — pure Go, no external packages
- **Compatibility layers** — drop-in adapters for survey/v2 and promptui
- **Functional options** — validators, transformers, defaults, help text
- **Structured forms** — `Ask()` fills struct fields from multiple prompts

## Install

```bash
go get github.com/agentine/inquire
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/agentine/inquire"
)

func main() {
	name, err := inquire.Input("What is your name?")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello, %s!\n", name)
}
```

## License

MIT
