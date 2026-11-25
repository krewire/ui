# `form` (`github.com/krewire/ui/form`)

## Purpose & Design Role

`form` provides a programmatic, server-rendered HTML form builder with request binding, field-level validation, and accessible markup generation. It enables Go handlers to declare forms, bind `http.Request` payloads, validate inputs, and render error-annotated HTML without hand-writing `<form>` tags or validation boilerplate.

All rendered labels, field names, placeholders, values, and validation messages are escaped via `html/template.HTMLEscapeString`, and every `<label for="field-<name>">` is paired with a matching `id="field-<name>"` on its control for accessibility.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C leaf package.
- **Dependencies**: Go standard library only (`fmt`, `html/template`, `net/http`, `net/mail`, `net/url`, `strconv`, `strings`). Enforced as a strict leaf module by `core/architecture_test.go`.
- **Dependents**: Consumed by `core` (`App.Form`).

## Core Concepts & APIs

### Form Container (`form.go`)
- `Form`: Holds `Action`, `Method` (defaults to `"POST"`), `Fields []Field`, `SubmitLabel` (defaults to `"Submit"`), `ValuesMap`, `ErrorsMap`, and `Class`.
- `New(action, method string) *Form`: Constructs a new form.
- Fluent builders: `Add(fields ...Field) *Form`, `WithSubmitLabel(label string) *Form`, `SetValue(name, val string) *Form`.
- Binding & validation:
  - `Bind(vals url.Values) *Form`: Populates `ValuesMap` from query or form values.
  - `BindRequest(r *http.Request) error`: Calls `r.ParseForm()` and binds `r.Form`.
  - `Validate() bool`: Runs each field's `Validate` method, populates `ErrorsMap`, and returns `true` if all fields pass.
  - `Value(name string) string`, `Values() map[string]string`, `Errors() map[string][]string`.
- `Render() template.HTML`: Emits `<form class="ui-form ...">` with all fields and a primary submit button.

### Field Contract & Concrete Fields
- `Field` interface (`field.go`): `Name() string`, `Label() string`, `Validate(value string) []string`, `Render(value string, errors []string) template.HTML`.
- `BaseField`: Shared struct implementing `Name()`, `Label()`, and `Validate()` (checks `Required` and runs `Validators []func(string) string`).
- `TextField` (`text_field.go`): `NewText(name, label string)`, `MakeRequired()`, `WithPlaceholder(ph string)`, `WithHelp(help string)`, `AsEmail()` (validates via `net/mail.ParseAddress`), `AsPassword()`, `MinLength(min int)`.
- `NumberField` (`number_field.go`): `NewNumber(name, label string)`, `SetMin(min float64)`, `SetMax(max float64)` (validates numeric parsing via `strconv.ParseFloat` and min/max bounds).
- `SelectField` (`select_field.go`): `NewSelect(name, label string, options ...SelectOption)`, `AddOption(value, label string)`. `SelectOption` holds `Value` and `Label`.
- `TextareaField` (`textarea_field.go`): `NewTextarea(name, label string)` (defaults to `Rows: 4`).
- `CheckboxField` (`checkbox_field.go`): `NewCheckbox(name, label string)` (renders inline control-before-label layout; treats `"true"`, `"on"`, and `"1"` as checked).

## Realistic Usage Examples

```go
package main

import (
	"fmt"
	"net/url"

	"github.com/krewire/ui/form"
)

func main() {
	f := form.New("/signup", "POST").
		WithSubmitLabel("Create Account").
		Add(
			form.NewText("email", "Email").MakeRequired().AsEmail().WithPlaceholder("you@example.com"),
			form.NewText("password", "Password").MakeRequired().AsPassword().MinLength(12),
			form.NewSelect("role", "Role",
				form.SelectOption{Value: "builder", Label: "Builder"},
				form.SelectOption{Value: "admin", Label: "Admin"},
			),
			form.NewCheckbox("terms", "I accept the terms"),
		)

	f.Bind(url.Values{
		"email":    {"ada@krewire.dev"},
		"password": {"short"},
		"role":     {"builder"},
		"terms":    {"on"},
	})

	if !f.Validate() {
		fmt.Println("Validation errors:", f.Errors())
	}
	html := f.Render()
	fmt.Printf("Rendered %d bytes of HTML\n", len(html))
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/form`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
