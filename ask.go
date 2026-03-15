package inquire

import (
	"fmt"
	"reflect"
)

// Question defines a single question in a structured form.
type Question struct {
	Name      string
	Prompt    Prompt
	Validate  func(any) error
	Transform func(any) any
}

// InputPrompt wraps Input as a Prompt.
type InputPrompt struct {
	Message string
	Options []Option
}

// Run executes the input prompt.
func (p *InputPrompt) Run() (any, error) {
	return Input(p.Message, p.Options...)
}

// PasswordPrompt wraps Password as a Prompt.
type PasswordPrompt struct {
	Message string
	Options []Option
}

// Run executes the password prompt.
func (p *PasswordPrompt) Run() (any, error) {
	return Password(p.Message, p.Options...)
}

// ConfirmPrompt wraps Confirm as a Prompt.
type ConfirmPrompt struct {
	Message string
	Options []Option
}

// Run executes the confirm prompt.
func (p *ConfirmPrompt) Run() (any, error) {
	return Confirm(p.Message, p.Options...)
}

// SelectPrompt wraps Select as a Prompt.
type SelectPrompt struct {
	Message string
	Items   []string
	Options []Option
}

// Run executes the select prompt. Returns the selected string.
func (p *SelectPrompt) Run() (any, error) {
	_, val, err := Select(p.Message, p.Items, p.Options...)
	return val, err
}

// MultiSelectPrompt wraps MultiSelect as a Prompt.
type MultiSelectPrompt struct {
	Message string
	Items   []string
	Options []Option
}

// Run executes the multi-select prompt. Returns []string.
func (p *MultiSelectPrompt) Run() (any, error) {
	_, vals, err := MultiSelect(p.Message, p.Items, p.Options...)
	return vals, err
}

// MultilinePrompt wraps Multiline as a Prompt.
type MultilinePrompt struct {
	Message string
	Options []Option
}

// Run executes the multiline prompt.
func (p *MultilinePrompt) Run() (any, error) {
	return Multiline(p.Message, p.Options...)
}

// EditorPrompt wraps Editor as a Prompt.
type EditorPrompt struct {
	Message string
	Options []Option
}

// Run executes the editor prompt.
func (p *EditorPrompt) Run() (any, error) {
	return Editor(p.Message, p.Options...)
}

// Ask runs a set of questions and fills the response struct by matching
// Question.Name to struct field names.
func Ask(questions []*Question, response any, opts ...Option) error {
	rv := reflect.ValueOf(response)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("inquire: response must be a pointer to a struct")
	}
	rv = rv.Elem()

	for _, q := range questions {
		answer, err := q.Prompt.Run()
		if err != nil {
			return err
		}

		if q.Validate != nil {
			if verr := q.Validate(answer); verr != nil {
				return verr
			}
		}

		if q.Transform != nil {
			answer = q.Transform(answer)
		}

		if q.Name == "" {
			continue
		}

		field := rv.FieldByName(q.Name)
		if !field.IsValid() {
			return fmt.Errorf("inquire: struct has no field %q", q.Name)
		}
		if !field.CanSet() {
			return fmt.Errorf("inquire: field %q is not settable (must be exported)", q.Name)
		}

		val := reflect.ValueOf(answer)
		if !val.Type().AssignableTo(field.Type()) {
			return fmt.Errorf("inquire: cannot assign %T to field %q of type %s", answer, q.Name, field.Type())
		}
		field.Set(val)
	}

	return nil
}
