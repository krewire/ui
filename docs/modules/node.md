# `node` (`github.com/krewire/ui/node`)

## Purpose & Design Role

`node` is the canonical, streaming HTML render-tree engine for programmatic Go UI. Every element, text node, layout wrapper, slot outlet, legacy widget adapter, and embedded `.kiw` component implements the `Node` interface, supporting both streaming output (`io.Writer`) and `html/template.HTML` return values.

Security and composability are built into the core primitives: text nodes are HTML-escaped by default, attribute keys are sanitized, URL attributes (`href`, `src`, `action`, `formaction`, `poster`, `data`) are validated through `SafeURL` to block `javascript:` / `vbscript:` / `data:` injection, and `context.Context` carries cancellation checks and named `SlotMap` overrides down the tree.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C leaf package.
- **Dependencies**: Go standard library only (`bytes`, `context`, `fmt`, `html/template`, `io`, `sort`, `strings`, `sync`, `unicode`). Enforced as a strict leaf module by `core/architecture_test.go`.
- **Dependents**: `core` (`RegisterNode` bridges `node.Node` into `ui.Registry`).

## Core Concepts & APIs

### Core Render Tree (`node.go`, `component.go`)
- `Node` interface:
  - `Render(ctx context.Context, w io.Writer) error`
  - `RenderHTML(ctx context.Context) (template.HTML, error)`
- `NodeFunc`: Function adapter implementing `Node`.
- `RenderToString(ctx context.Context, n Node) (string, error)` / `RenderToHTML(ctx context.Context, n Node) (template.HTML, error)`.
- Node constructors:
  - `Text(s string) Node`, `Textf(format string, a ...any) Node`, `Strings(strs ...string) []Node`: Contextually HTML-escaped text nodes.
  - `RawHTML(h string) Node`: Unescaped output for trusted SVG or pre-sanitized HTML.
  - `Fragment(children ...Node) Node`: Groups nodes without a wrapper element.
  - `Element(tag string, attrs Attrs, children ...Node) Node`: Generic HTML element supporting void tags (`img`, `input`, `br`, `hr`, etc.), sorted deterministic attributes (`Attrs`), and `SafeURL` filtering.
  - `SafeURL(raw string) string`: Allows relative URLs and `http`/`https`/`mailto`/`tel`/`ftp` schemes; replaces unsafe schemes with `"#"`.
- `ComponentFunc[P any]`: Generic functional component signature `func(props P, children ...Node) Node`.
- `HTMLRenderer` & `FromWidget(w HTMLRenderer) Node`: Adapts any `Render() template.HTML` type (such as `widget.Widget`) into a `Node`.

### Layout & UI Primitives (`layout.go`, `primitive.go`, `props.go`)
- Typed prop enums: `Variant` (`VariantDefault`, `VariantPrimary`, `VariantSecondary`, `VariantOutline`, `VariantGhost`, `VariantError`, `VariantSuccess`, `VariantWarning`, `VariantInfo`), `Size` (`SizeSmall`, `SizeMedium`, `SizeLarge`), `Direction`, `Align`, `Justify`, `HeadingLevel`, `InputType`, plus `ClassNames(classes ...string) string`.
- Layout nodes: `Container(ContainerProps, ...Node)`, `Stack(StackProps, ...Node)`, `VStack`, `HStack`, `Flex`, `Grid(GridProps, ...Node)`, `Section(SectionProps, ...Node)`, `Box(BoxProps, ...Node)`.
- UI primitives: `Button(ButtonProps, ...Node)`, `Input(InputProps)`, `Heading(HeadingProps, ...Node)`, `Paragraph(TextProps, ...Node)`, `Badge(BadgeProps, ...Node)`, `Card(CardProps, ...Node)`, `Alert(AlertProps, ...Node)`, `Divider(DividerProps)`, `Modal(ModalProps, ...Node)`, `Table(TableProps)`.

### Named Slots & Go ↔ `.kiw` Interop (`slot.go`, `kiw_adapter.go`)
- `SlotMap` (`map[string]Node`), `Slot(name string, fallback ...Node) Node`, `WithSlots(child Node, slots SlotMap) Node`, `ContextWithSlots`, `SlotsFromContext`: Context-propagated named slot outlets with fallback content.
- `Kiw(name string, props map[string]any, children ...Node) Node`: Embeds a `.kiw` component inside a Go `Node` tree, resolved via `RegisterKiwComponent` or `SetKiwFallbackResolver` (falling back gracefully to a `data-kiw-component` placeholder if unregistered).

## Realistic Usage Examples

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/krewire/ui/node"
)

func main() {
	page := node.Container(node.ContainerProps{MaxWidth: "max-w-4xl", Centered: true},
		node.VStack(node.StackProps{Gap: "gap-4"},
			node.Heading(node.HeadingProps{Level: node.Heading1}, node.Text("Deployments")),
			node.Card(node.CardProps{
				Title: "Production Cluster",
				Slots: node.SlotMap{
					"footer": node.Button(node.ButtonProps{Variant: node.VariantPrimary}, node.Text("Restart")),
				},
			}, node.Paragraph(node.TextProps{Muted: true}, node.Text("All replicas healthy."))),
		),
	)

	html, err := node.RenderToString(context.Background(), page)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html)
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/node`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
