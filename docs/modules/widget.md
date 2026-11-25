# `widget` (`github.com/krewire/ui/widget`)

## Purpose & Design Role

`widget` provides a fluent, builder-style Go API for rendering standalone HTML UI components (`Button`, `Badge`, `Card`, `Alert`, `Modal`, `Table`, `Grid`, `Stack`, `Heading`, `Text`, `Divider`). Every widget implements the `Widget` interface (`Render() template.HTML`), escapes user-supplied strings via `html/template.HTMLEscapeString`, and sanitizes link targets with `SafeURL`.

Styling uses `--ui-*` CSS custom properties and `ui-*` class prefixes (`tokens.go`, `style.go`) so widgets integrate seamlessly with `styles.CSS` and `ui.Theme`.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C leaf package.
- **Dependencies**: Go standard library only (`fmt`, `html/template`, `strconv`, `strings`). Enforced as a strict leaf module by `core/architecture_test.go`.
- **Dependents**: `panel` and `core` (`App`).

## Core Concepts & APIs

### Core Contract & Security (`widget.go`, `fragment.go`)
- `Widget` interface: `Render() template.HTML`.
- `HTML` (`type HTML string`): Adapts trusted raw HTML into a `Widget`.
- `Fragment(widgets ...Widget) template.HTML`: Concatenates multiple rendered widgets without a wrapper element.
- `SafeURL(raw string) string`: Strips control characters, permits relative URLs and `http`/`https`/`mailto`/`tel`/`ftp` schemes, and replaces unsafe schemes (`javascript:`, `vbscript:`, `data:`) with `"#"`.

### Design Tokens & Layout Bounds (`tokens.go`, `bounds.go`)
- Variants: `VariantPrimary`, `VariantSecondary`, `VariantSuccess`, `VariantWarning`, `VariantError`, `VariantInfo`, `VariantDefault`.
- Sizes: `SizeSmall` (`"sm"`), `SizeMedium` (`"md"`), `SizeLarge` (`"lg"`).
- Directions: `DirectionVertical`, `DirectionHorizontal`.
- Class prefixes: `ButtonClassPrefix` (`"ui-btn"`), `BadgeClassPrefix` (`"ui-badge"`), `AlertClassPrefix`, `CardClassPrefix`, `TableClassPrefix`, `ModalClassPrefix`, `StackClassPrefix`, `GridClassPrefix`, `HeadingClassPrefix`.
- Layout constraints (`bounds.go`): Heading levels clamped to `1..6` (default `2`), `Grid` columns clamped to `2..4` (default `2`).

### Component Builders
- `Button` (`NewButton(text string)`): `Primary()`, `Secondary()`, `Error()`, `Small()`, `Large()`, `Link(url string)`, `WithIcon(icon string)`, `Submit()`, `Disable()`, `WithClickHandler(handler string)`, `WithClass(class string)`.
- `Badge` (`NewBadge(text string)`): `Primary()`, `Secondary()`, `Success()`, `Warning()`, `Error()`, `WithIcon(icon string)`.
- `Alert` (`NewAlert(message string)`): `Success()`, `Warning()`, `Error()` (pairs variants with semantic icons).
- `Card` (`NewCard(title string)`): `WithDesc`, `WithBody`, `WithFooter`, `WithActions`.
- `Modal` (`NewModal(id, title string)`): `WithContent`, `WithFooter`.
- `Table` (`NewTable(headers ...string)`): `AddRow(cells ...template.HTML)`, `AddTextRow(cells ...string)`.
- `Stack` (`VStack(widgets ...Widget)`, `HStack(widgets ...Widget)`): `WithGap(gap string)`.
- `Grid` (`NewGrid(columns int, widgets ...Widget)`), `Heading` (`NewHeading(level int, content string)`), `Text` (`NewText(content string)` with `Strong()`, `Subtle()`), `Divider` (`NewDivider(label ...string)`).

## Realistic Usage Examples

```go
package main

import (
	"fmt"

	"github.com/krewire/ui/widget"
)

func main() {
	tbl := widget.NewTable("Service", "Status", "Region").
		AddTextRow("api-edge", "Healthy", "us-east").
		AddTextRow("worker-01", "Degraded", "eu-west")

	card := widget.NewCard("Cluster Inventory").
		WithDesc("Live node status").
		WithActions(widget.NewBadge("Live").Success()).
		WithBody(tbl).
		WithFooter(
			widget.HStack(
				widget.NewButton("Export CSV").Secondary().Small(),
				widget.NewButton("Deploy").Primary().Small(),
			).WithGap("0.5rem"),
		)

	fmt.Println(card.Render())
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/widget`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
