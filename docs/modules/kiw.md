# `kiw` (`github.com/krewire/ui/kiw`)

## Purpose & Design Role

`kiw` is the canonical Go-native parser, template desugarer, linter, diagnostic mapper, and package bundler for the `.kiw` single-file component format. A `.kiw` file combines optional YAML frontmatter (`---`), Svelte-inspired template markup (`{#if}`, `{#each}`, `{#snippet}`, `{@render}`, `<slot>`), JSX-like PascalCase component tags (`<Card title={title} />`), scoped `<style>` blocks, `<script>` island blocks, and `<markdown>` sections.

Because `kiw` is a pure-Go leaf package, static-site generators (`web/ssg`), UI libraries, and the `krew` CLI can all parse and validate `.kiw` files without any import cycles or Node.js runtime.

## Architectural Position

- **Layer**: Layer 2 (Modular Domain Libraries & Design System).
- **Distribution Model**: Model C leaf package.
- **Dependencies**: Go standard library only (`archive/tar`, `bytes`, `compress/gzip`, `crypto/sha256`, `encoding/json`, `errors`, `fmt`, `io`, `os`, `path/filepath`, `regexp`, `sort`, `strconv`, `strings`, `sync`, `text/template`, `unicode`, `unicode/utf8`). Enforced as a strict leaf module by `core/architecture_test.go`.
- **Dependents**: `dsl` (aliases `kiw`), `editors` (runs `kiw.Lint` and `kiw.ParseKiw` in tests), `core`, and `web/ssg`.

## Core Concepts & APIs

### Parsing & Module AST (`kiw.go`, `imports.go`)
- `KiwModule`: Parsed result containing `Frontmatter`, `Body` (desugared Go `html/template` source), `Styles`, `Scripts`, `StyleBlocks []StyleBlock`, `ScriptBlocks []ScriptBlock`, `Markdown`, `Components []ComponentCall`, `ComponentNames`, `Imports []ImportDecl`, `Exports []ExportDecl`, `Props map[string]PropSpec`, `Defaults`, and `Slots []SlotSpec`.
- `ParseKiw(src string) (*KiwModule, error)` / `ParseKiwFile(path string) (*KiwModule, error)`: Extracts frontmatter, `<style>`, `<script>`, `<markdown>`, imports/exports/props/slots, and desugars the template body.
- `MarkdownRenderer func(src []byte) (string, error)`: Optional pluggable hook for rendering `<markdown>` blocks.
- `GetBuffer()` / `PutBuffer(buf *bytes.Buffer)`: Package-level `sync.Pool` buffer helpers.

### Template Desugaring & Prop Resolution (`desugar.go`, `desugar_comp.go`, `desugar_expr.go`, `props.go`)
- `DesugarTemplate(src string) (string, error)` / `DesugarTemplateWithScope(src string, loopVars map[string]bool) (string, error)`:
  - Translates `class:active={isActive}` directives into conditional `{{if ...}}` class attributes.
  - Translates `{#each items as item, index}...{/each}` into Go `{{range}}` blocks with nested scope tracking.
  - Translates `{#if}`, `{:else if}`, `{:else}`, `{/if}` into Go template conditionals.
  - Translates `{#snippet name(params)}...{/snippet}` and `{@render name(args)}` into `{{define}}` / `{{template}}` calls.
  - Translates `<ComponentName prop={val} />` and paired `<Card>...</Card>` into `{{component ...}}` / `{{mount ...}}` calls.
  - Translates `<slot>` / `<slot name="header">fallback</slot>` and `{var.field}` interpolations into `{{.Body}}` / `{{.Var.Field}}`.
- `PropAliases(name string) []string`: Implements lenient prop-casing (`variant` <-> `Variant`, `aria-label` -> `ariaLabel` / `AriaLabel`, initialisms like `id` -> `ID`, `url` -> `URL`).

### Linting & Line-Accurate Diagnostics (`lint.go`, `diag.go`)
- `Diagnostic`: Error struct with `File`, 1-based `Line`/`Col`, `Message`, `Hint`, and `Err` (`Unwrap() error`).
- `Lint(file, src string) []Diagnostic`: Catches invalid YAML frontmatter, control blocks inside component tag attributes, and unbalanced `{#if}`/`{#each}`/`{#snippet}` blocks.
- `AlignedBody(src string) string`: Produces a line-preserving desugared body where stripped blocks become blank lines so template line numbers match `.kiw` source lines.
- `Explain(file, name, src string, funcs map[string]any, err error) error`: Maps `text/template` or `html/template` parse/exec errors back to a `*Diagnostic` positioned in the original `.kiw` source with actionable hints.

### Standalone Component Packages (`pkg.go`)
- Constants: `PackageManifestFile` (`"kiw.pkg.json"`), `PackageExtension` (`".kpkg"`), `DefaultPackageDir` (`".krewire/pkg"`).
- `PackageManifest`: Defines `Name`, `Version`, `Components` (map of component name to relative `.kiw` path), `Styles`, `Scripts`, and `Dependencies`.
- `LoadPackageManifest`, `SavePackageManifest`, `ValidatePackageManifest` (enforces `SafePath` relative-path bounds), and `PackPackage(dir, outPath string)` (produces a `.kpkg` tar.gz archive and SHA-256 hex digest).

## Realistic Usage Examples

```go
package main

import (
	"fmt"
	"log"

	"github.com/krewire/ui/kiw"
)

func main() {
	src := `---
title: Dashboard
---
{#if showBanner}
  <Alert variant="info">{message}</Alert>
{/if}`

	if diags := kiw.Lint("page.kiw", src); len(diags) > 0 {
		for _, d := range diags {
			log.Println(d.Error())
		}
		return
	}

	mod, err := kiw.ParseKiw(src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Components used:", mod.ComponentNames)
	fmt.Println("Desugared body:\n", mod.Body)
}
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/kiw`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
