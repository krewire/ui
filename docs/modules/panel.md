# `panel` (`github.com/krewire/ui/panel`)

## Purpose & Design Role

`panel` provides higher-level administrative UI containers—`Panel`, `Stat`, and `Dashboard`—built on top of `widget.Widget`. It is designed for operational consoles, developer dashboards, and admin views that need structured cards with header actions, KPI metric tiles with trend indicators, and responsive multi-column dashboard layouts.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C package distribution.
- **Dependencies**: Imports `github.com/krewire/ui/widget` and the Go standard library (`fmt`, `html/template`, `strings`). Enforced by `core/architecture_test.go` to never import `core`, `form`, or `kiw`.
- **Dependents**: Consumed by `core` (`App.Panel` and `App.Dashboard`).

## Core Concepts & APIs

All three types in `panel/panel.go` implement `widget.Widget` (`Render() template.HTML`):

### `Panel`
- `Panel`: Holds `Title`, `Description`, `Actions widget.Widget`, `Body widget.Widget`, and `Footer widget.Widget`.
- `New(title string) *Panel`: Creates a new panel.
- Fluent builders: `WithDesc(desc string) *Panel`, `WithActions(actions widget.Widget) *Panel`, `WithBody(body widget.Widget) *Panel`, `WithFooter(footer widget.Widget) *Panel`.
- `Render() template.HTML`: Emits `<div class="ui-panel">` with an escaped title/description header, optional action slot, body, and footer.

### `Stat`
- `Stat`: KPI metric card with `Label`, `Value`, `Change`, `Up bool`, and `Icon`.
- `NewStat(label, value string) *Stat`: Creates a stat card.
- Fluent builders: `WithChange(change string, up bool) *Stat` (renders `↑` with `ui-stat-change up` or `↓` with `ui-stat-change down`), `WithIcon(icon string) *Stat`.
- `Render() template.HTML`: Emits `<div class="ui-stat">`.

### `Dashboard`
- `Dashboard`: Organizes `Title`, `Stats []*Stat`, `Panels []*Panel`, and `Actions widget.Widget`.
- `NewDashboard(title string) *Dashboard`: Creates a dashboard.
- Fluent builders: `AddStat(stats ...*Stat) *Dashboard`, `AddPanel(panels ...*Panel) *Dashboard`, `WithActions(actions widget.Widget) *Dashboard`.
- `Render() template.HTML`: Emits `<div class="ui-dashboard">` with a top header, a responsive stats row (`ui-grid-1` through `ui-grid-4`, capped at 4 columns), and stacked panels.

## Realistic Usage Examples

```go
package main

import (
	"fmt"

	"github.com/krewire/ui/panel"
	"github.com/krewire/ui/widget"
)

func main() {
	dash := panel.NewDashboard("Service Telemetry").
		WithActions(widget.NewButton("Refresh").Primary()).
		AddStat(
			panel.NewStat("Requests/sec", "1,420").WithChange("+12%", true),
			panel.NewStat("Error Rate", "0.02%").WithChange("-0.01%", false),
		).
		AddPanel(
			panel.New("Recent Events").
				WithDesc("Last 24 hours").
				WithBody(widget.NewAlert("Zero incidents recorded").Success()),
		)

	html := dash.Render()
	fmt.Printf("Dashboard HTML length: %d\n", len(html))
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/panel`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
