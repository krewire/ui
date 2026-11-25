# Krewire UI (`github.com/krewire/ui`) — Module Reference

Krewire UI is a pure-Go, CGO-free multi-module monorepo organized around **Model C package distribution** ([ADR-K7P9L](https://github.com/krewire/internal)) and the **6-Layer Pyramid Architecture** ([ADR-E8P9D](https://github.com/krewire/internal)). There is no root `go.mod`; instead, `go.work` coordinates 10 independent Go modules.

Nine of the ten modules are self-contained leaf or near-leaf packages with zero external dependencies, while `core` (`package ui`) serves as the ergonomic composition facade.

## Module Catalog

### Core & Composition
- [`core`](core.md) (`github.com/krewire/ui/core`, package `ui`) — Umbrella facade providing `App`, `Layout`, `Theme`, `Palette`, `Registry`, declarative `LayoutFromConfig`, and CSS helpers.

### Template Engine & Component Catalog
- [`components`](components.md) (`github.com/krewire/ui/components`) — Embedded catalog of 66+ clean, modern `.kiw` single-file components (`FS`, `Registry`, `Resolve`, starter `Packs`).
- [`kiw`](kiw.md) (`github.com/krewire/ui/kiw`) — Canonical `.kiw` parser, Svelte-style template desugarer (`{#if}`, `{#each}`, `{#snippet}`, `<slot>`), linter, diagnostic explainer, and `.kpkg` packager.
- [`dsl`](dsl.md) (`github.com/krewire/ui/dsl`) — Backward-compatible forwarding alias for `github.com/krewire/ui/kiw`.
- [`editors`](editors.md) (`github.com/krewire/ui/editors`) — TextMate grammar (`text.html.kiw`), VS Code extension manifest, language configuration, and grammar validation tests.

### Programmatic Go UI Primitives
- [`node`](node.md) (`github.com/krewire/ui/node`) — Streaming `Node` render tree (`io.Writer` + `template.HTML`), `Element`, `SafeURL` XSS guard, named `Slot` context propagation, layout primitives, and Go ↔ `.kiw` `Kiw()` adapter.
- [`widget`](widget.md) (`github.com/krewire/ui/widget`) — Fluent, builder-style `Widget` primitives (`Button`, `Badge`, `Card`, `Alert`, `Modal`, `Table`, `Grid`, `Stack`, `Heading`, `Text`, `Divider`) with `SafeURL` sanitization.
- [`form`](form.md) (`github.com/krewire/ui/form`) — Programmatic `Form` builder with `http.Request` binding, field-level validators, and accessible HTML rendering (`TextField`, `NumberField`, `SelectField`, `TextareaField`, `CheckboxField`).
- [`panel`](panel.md) (`github.com/krewire/ui/panel`) — Higher-level administrative `Panel`, KPI `Stat`, and `Dashboard` widgets built on `widget.Widget`.

### Design System & Tokens
- [`styles`](styles.md) (`github.com/krewire/ui/styles`) — Embedded canonical clean, modern stylesheet (`krewire.css` via `styles.CSS`) with `--ui-*` CSS custom properties, typography/spacing scales, dark mode, and layout utilities.

## Dependency & Layering Rules

Enforced by `core/architecture_test.go`:
1. **Strict Leaf Packages**: `node`, `widget`, `form`, `styles`, and `kiw` must never import `github.com/krewire/ui` or any sibling `ui/*` module.
2. **Scoped Intermediate Packages**: `panel` imports only `widget` (never `core`, `form`, or `kiw`); `dsl` imports only `kiw`.
3. **Ecosystem Isolation**: No module in `ui` may import Layer 3+ repositories (`krew`, `cloud`, `auth`, `dash`).
4. **Pure Go (`CGO_ENABLED=0`)**: Every module compiles with the Go standard library only.
