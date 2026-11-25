# `styles` (`github.com/krewire/ui/styles`)

## Purpose & Design Role

`styles` is the canonical CSS design-system distribution module for Krewire UI. It embeds `krewire.css`—containing the clean, modern design tokens, CSS custom properties (`--ui-*`), typography and spacing scales, elevation/pop-shadow tokens, dark-mode variables, base resets, atomic control styles, and responsive layout primitives—into a single Go string constant.

Shipping the stylesheet as an embedded Go package means applications and static-site generators can serve or inline the exact canonical CSS at runtime with zero external asset pipelines or Node.js build steps.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C leaf package (`go:embed krewire.css`).
- **Dependencies**: Go standard library only (`embed`). Zero `require` lines in `styles/go.mod`. Enforced as a strict leaf module by `core/architecture_test.go`.
- **Dependents**: `core` and downstream web servers/SSG pipelines that serve `/assets/krewire.css` or inline `styles.CSS`.

## Core Concepts & APIs

- `CSS` (`string` in `styles/styles.go`): Embedded contents of `styles/krewire.css`.
- Token groups defined inside `krewire.css`:
  - **Typography**: `--ui-font-sans`, `--ui-font-display`, `--ui-font-mono`, `--ui-text-2xs` through `--ui-text-5xl`, line-height (`--ui-leading-*`), letter-spacing (`--ui-tracking-*`), and font-weight (`--ui-font-*`) tokens.
  - **Spacing & Geometry**: 4px/8px grid tokens (`--ui-space-0` through `--ui-space-24`), border radii, and neo-brutalist pop borders/shadows.
  - **Colors & Dark Mode**: Semantic surface, text, border, and accent variables supporting light and `[data-theme="dark"]` modes.
  - **Atomic Controls & Layout Primitives**: Styling for `.ui-*` widgets, forms, panels, and `.kiw-*` components.

## Realistic Usage Examples

Serving the canonical stylesheet directly from a standard Go `http.ServeMux`:

```go
package main

import (
	"log"
	"net/http"

	"github.com/krewire/ui/styles"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/assets/krewire.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write([]byte(styles.CSS))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/styles`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
