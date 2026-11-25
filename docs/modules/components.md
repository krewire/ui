# `components` (`github.com/krewire/ui/components`)

## Purpose & Design Role

`components` ships the canonical Krewire UI (kui) catalog of reusable `.kiw` single-file components as an embeddable Go module. Rather than forcing applications to copy-paste template files into their repository, `components` embeds all `*.kiw` component sources alongside a machine-verified `components.json` manifest.

Downstream consumers—including `core` (`ui.ComponentsFS`), static-site generators (`web/ssg`), and the `krew` CLI (`krew ui vendor` / `krew add ui:<pack>`)—read component sources and dependency metadata directly from this package.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C package distribution (`go:embed *.kiw components.json`).
- **Dependencies**: Go standard library only (`embed`, `encoding/json`, `fmt`, `sort`, `strings`). Zero external `require` directives in `components/go.mod`.
- **Dependents**: `github.com/krewire/ui/core` re-exports `components.FS` as `ui.ComponentsFS`; `web/ssg` loads built-in `.kiw` tags from it.

## Core Concepts & APIs

### Embedded Filesystem & Registry
- `FS` (`embed.FS`): Embeds every `*.kiw` file at the filesystem root together with `components.json`.
- `ManifestFile` (`"components.json"`): The embedded catalog file, validated against `components.schema.json` (JSON Schema Draft-07).
- `Registry` (`map[string]ComponentMeta`): Package-level map initialized at startup from `components.json` (panics via `mustRegistry()` if the embedded manifest is malformed).
- `Version() string`: Returns the component catalog format version (`"0.2.0"`).

### Discovery & Lookup
- `List() []string`: Returns a sorted slice of all component names in `Registry`.
- `Get(name string) ([]byte, error)`: Reads raw `.kiw` template source for `name` (accepts `"Button"` or `"Button.kiw"`). Wraps `fs.ErrNotExist` when missing.
- `Lookup(name string) (ComponentMeta, bool)`: Case-insensitive lookup (`"button"`, `"Button"`, `"Button.kiw"`) returning the `ComponentMeta` entry.

### Manifest, Packs & Dependency Resolution
- `ComponentMeta`: Describes a single component (`Name`, `Category`, `Description`, `File`, `Props`, `Defaults`, `Dependencies`, `Styles`, `Script`, `Preview`). `Props`, `Defaults`, `Dependencies`, `Styles`, and `Script` are derived from the `.kiw` markup by `go generate ./components` (`internal/manifestgen`).
- `Pack`: A curated starter bundle (`Name`, `Description`, `Components`) installed via `krew add ui:<name>`.
- `Manifest`: Root decoded structure (`Schema`, `Version`, `Packs`, `Components`).
- `LoadManifest() (*Manifest, error)`: Decodes `components.json` with `DisallowUnknownFields()` enabled.
- `Packs() []Pack`: Returns all starter packs sorted by name.
- `Resolve(refs ...string) ([]string, error)`: Expands component names and `"pack:<name>"` references into a sorted, deduplicated list of components including all transitive `Dependencies`.

## Realistic Usage Examples

```go
package main

import (
	"fmt"
	"log"

	"github.com/krewire/ui/components"
)

func main() {
	// List all available components
	names := components.List()
	fmt.Printf("Catalog contains %d components\n", len(names))

	// Case-insensitive metadata lookup
	if meta, ok := components.Lookup("button"); ok {
		fmt.Printf("%s (%s): props=%v\n", meta.Name, meta.Category, meta.Props)
	}

	// Resolve a starter pack plus an extra component with transitive deps
	toInstall, err := components.Resolve("pack:landing", "Modal")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Resolved set:", toInstall)

	// Read raw .kiw source
	src, err := components.Get("Button")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(src[:min(len(src), 80)]))
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/components` (or repository root via `go.work`):

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
