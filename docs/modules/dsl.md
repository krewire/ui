# `dsl` (`github.com/krewire/ui/dsl`)

## Purpose & Design Role

`dsl` is a thin, backward-compatible forwarding alias for [`github.com/krewire/ui/kiw`](kiw.md). Historically, the `.kiw` template parser lived in `dsl`; when static-site generation (`web/ssg`) and lower-level packages needed to parse `.kiw` files without introducing an import cycle with higher-level tooling, the canonical implementation moved to `github.com/krewire/ui/kiw`.

`dsl` preserves existing import paths via Go type aliases and wrapper functions while directing all new code to import `github.com/krewire/ui/kiw` directly. The directory also contains `kiw.ts`, a TypeScript reference helper for JS/TS consumers.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C package distribution.
- **Dependencies**: Imports only `github.com/krewire/ui/kiw` (verified by `core/architecture_test.go`, which forbids `dsl` from importing `core`, `form`, or `panel`).

## Core Concepts & APIs

All exported symbols in `dsl/kiw.go` are direct aliases or pass-through wrappers to `kiw`:

- `type KiwModule = kiw.KiwModule` — Parsed representation of a `.kiw` file (frontmatter, body, styles, scripts, markdown, components, imports, exports, props, defaults, slots).
- `type StyleBlock = kiw.StyleBlock` — Extracted `<style>` block with `Lang`, `Scoped`, and `Content`.
- `type ScriptBlock = kiw.ScriptBlock` — Extracted `<script>` block with `Lang`, `Hydrate`, `Server`, `Compute`, and `Content`.
- `type ComponentCall = kiw.ComponentCall` — Recorded component invocation (`Name`, `Hydrate`, `Framework`, `Props`, `Raw`).
- `DesugarTemplate(src string) (string, error)` — Translates `.kiw` syntax (`{#if}`, `{#each}`, `<Component />`, `{prop}`) into Go `html/template` syntax.
- `ParseKiw(src string) (*KiwModule, error)` — Parses `.kiw` source into a `*KiwModule`.
- `ParseKiwFile(path string) (*KiwModule, error)` — Reads and parses a `.kiw` file from disk.

## Realistic Usage Examples

For new code, prefer importing `github.com/krewire/ui/kiw` directly. Existing callers using `dsl` continue to work identically:

```go
package main

import (
	"fmt"
	"log"

	"github.com/krewire/ui/dsl"
)

func main() {
	src := "---\ntitle: Welcome\n---\n<Card title={title}>Hello</Card>"
	mod, err := dsl.ParseKiw(src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Title: %v, Components: %v\n", mod.Frontmatter["title"], mod.ComponentNames)
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/dsl`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
