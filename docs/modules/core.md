# `core` (`github.com/krewire/ui/core`, package `ui`)

## Purpose & Design Role

`core` is the primary public facade of the Krewire UI framework. While the other nine modules provide focused, decoupled leaf capabilities (DOM nodes, widgets, forms, panels, CSS tokens, `.kiw` parsing), `core` composes them into an ergonomic application-building layer: an HTTP `App` builder, a `Layout` page shell, a light/dark `Theme` and `Palette` system, declarative `LayoutFromConfig` assembly, and a typed `Registry` for programmatic Go components.

## Architectural Position

- **Layer**: Layer 2 (Composition Root within `github.com/krewire/ui`).
- **Distribution Model**: Model C umbrella package (`package ui`).
- **Dependencies**: Imports all nine sibling modules (`components`, `dsl`, `editors`, `form`, `kiw`, `node`, `panel`, `styles`, `widget`) plus the Go standard library (`html/template`, `net/http`, `strings`).
- **Invariants**: Enforced by `architecture_test.go`:
  - `TestArchitecture_SubsystemLayering` ensures leaf modules never import `core`.
  - `TestArchitecture_EcosystemIsolation` ensures no module in `ui` imports `krew`, `cloud`, `auth`, or `dash`.

## Core Concepts & APIs

### Programmatic Application Builder (`app_builder.go`)
- `App`: Implements `http.Handler` via `ServeHTTP` with an internal `http.ServeMux` and route tracking.
- `New(title string) *App`: Creates an `App` initialized with `DefaultTheme()` and `DefaultTailwindCDN` (`"https://cdn.tailwindcss.com"`).
- Fluent configuration: `WithTheme(t *Theme)`, `WithTailwindURL(url string)`, `DisableTailwind()`, `AddNav(label, url string)`.
- Route registration:
  - `Page(path, title string, content widget.Widget) *App`
  - `Panel(path string, p *panel.Panel) *App`
  - `Dashboard(path string, d *panel.Dashboard) *App`
  - `Form(path string, f *form.Form, onSubmit func(values map[string]string) (widget.Widget, error)) *App` (handles both `GET` render and `POST` bind/validate/submit).
- `RenderShell(pageTitle string, content template.HTML) template.HTML`: Wraps page content in a complete HTML5 document with header, nav, theme toggle, and Tailwind config.
- `ComponentsFS`: Re-exports `components.FS`.

### Page Shell & Declarative Layout (`layout.go`, `layout_config.go`)
- `Layout`: Composable HTML5 page shell with `Title`, `Description`, `Head`, `BodyAttrs`, `Header`, `Main`, `Footer`, `Aside`, `Stylesheets`, and `Theme`. Fluent helpers: `WithHead`, `WithBodyAttrs`, `WithTheme`.
- `Slot(parts ...template.HTML) template.HTML`: Concatenates pre-rendered HTML fragments into a layout region.
- `LayoutConfig`, `HeaderConfig`, `FooterConfig`, `NavConfig`, `NavLinkConfig`, `ComponentRef`, `ThemeConfig`: JSON/YAML-serializable configuration structs.
- `LayoutFromConfig(config *LayoutConfig) *Layout`: Constructs a `Layout` from declarative configuration.

### Theme & Palette System (`theme.go`, `palette.go`, `theme_toggle.go`)
- `Palette`: Defines semantic `Color` tokens (`Primary`, `PrimaryContent`, `Secondary`, `SecondaryContent`, `Accent`, `Base1`, `Base2`, `NeutralContent`, etc.) with `CSSVars(defaults Palette) string`.
- `DefaultLightPalette` / `DefaultDarkPalette`: Canonical light and dark clean, modern color palettes.
- `Theme`: Manages light/dark `Palette`s, `StorageKey` (defaults to `"krewire-theme"`), and `Default` (`"auto"`, `"light"`, `"dark"`).
  - `DefaultTheme() *Theme`, `WithStorageKey`, `WithDefault`
  - `Style() template.CSS`: Emits `:root`, `:root[data-theme="dark"]`, and `@media (prefers-color-scheme: dark)` CSS custom properties.
  - `Script() template.HTML`: Emits palette styles plus an inline pre-paint theme script (preventing FOUC) wired to `[data-theme-toggle]`.
  - `Button() template.HTML`: Emits an accessible sun/moon toggle button.
  - `TailwindConfigScript() template.HTML`: Maps `--ui-*` tokens into Tailwind utility colors.

### Component Registry & Built-in Renderers (`registry.go`, `components.go`, `components_css.go`)
- `Component` interface (`Render(props any) (template.HTML, error)`), `ComponentFunc`, `ComponentMeta`, `PropSchema`.
- `Registry`: Thread-safe component registry (`NewRegistry()`, `Default()`). Methods: `Register`, `RegisterWithMeta`, `Get`, `Has`, `Meta`, `AllMeta`, `Names`, `Render`.
- `RegisterNode[P any](r *Registry, name string, fn func(props P, children ...node.Node) node.Node, meta ...ComponentMeta) *Registry`: Generic bridge registering a `node.Node`-based functional component.
- Built-in Go structs (`Header`, `Footer`, `Nav`, `NavLink`, `Hero`, `Card`, `Button`, `Badge`, `Container`, `Grid`) and `AllBuiltinNames() []string`.
- `ComponentsCSS() string`: Returns the bundled layout/component CSS (`ContainerCSS`, `GridCSS`, `HeaderCSS`, `FooterCSS`, `CardCSS`, `HeroCSS`, `BadgeCSS`, `ButtonCSS`, etc.).

## Realistic Usage Examples

```go
package main

import (
	"log"
	"net/http"

	ui "github.com/krewire/ui/core"
	"github.com/krewire/ui/form"
	"github.com/krewire/ui/panel"
	"github.com/krewire/ui/widget"
)

func main() {
	app := ui.New("Control Plane").
		AddNav("Overview", "/").
		AddNav("Settings", "/settings")

	dash := panel.NewDashboard("System Overview").
		AddStat(panel.NewStat("Active Nodes", "24").WithChange("+3", true)).
		AddPanel(panel.New("Status").WithBody(widget.NewAlert("All systems nominal").Success()))

	settingsForm := form.New("/settings", "POST").
		Add(form.NewText("site_name", "Site Name").MakeRequired()).
		WithSubmitLabel("Save")

	app.Dashboard("/", dash)
	app.Form("/settings", settingsForm, func(vals map[string]string) (widget.Widget, error) {
		return widget.NewAlert("Saved settings for " + vals["site_name"]).Success(), nil
	})

	log.Fatal(http.ListenAndServe(":8080", app))
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/core`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
