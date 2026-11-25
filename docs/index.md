# Krewire UI — Developer & Contributor Documentation

Welcome to the internal engineering documentation for **`github.com/krewire/ui`**, the pop-brutalist `.kiw` design system, component catalog, and layout engine for the Krewire digital SDLC ecosystem.

## Documentation Index

- [Architecture & Modular Subsystems](architecture.md): 6-Layer pyramid integration, Model C package distribution, subsystem boundaries (`node`, `widget`, `form`, `styles`, `panel`, `kiw`), and layering invariants.
- [Design System & Tokens](design-system.md): Pop-Brutalist design tokens, typography scales, spacing grids, high-contrast surfaces, and dark mode palettes.
- [Getting Started](getting-started.md): Installation, package-based imports, vendoring components (`krew ui vendor`), layouts, and page integration.
- [Module Reference](modules/index.md): Comprehensive technical reference for all 10 Go modules in `github.com/krewire/ui`.
- [Contributing Guide](../CONTRIBUTING.md): Quality gates, adding components, tests, and contribution standards.

## Module Reference

- [Modules Overview](modules/index.md): Architecture, dependency rules, and module catalog.
- [`core`](modules/core.md) (`github.com/krewire/ui/core`): Public facade, `App` builder, `Layout`, `Theme`, `Palette`, and `Registry`.
- [`components`](modules/components.md) (`github.com/krewire/ui/components`): Embedded `.kiw` component catalog (`FS`, `Registry`, `Resolve`, starter `Packs`).
- [`kiw`](modules/kiw.md) (`github.com/krewire/ui/kiw`): `.kiw` parser, template desugarer, linter, diagnostic explainer, and `.kpkg` packager.
- [`dsl`](modules/dsl.md) (`github.com/krewire/ui/dsl`): Backward-compatible forwarding alias for `github.com/krewire/ui/kiw`.
- [`editors`](modules/editors.md) (`github.com/krewire/ui/editors`): TextMate grammar (`text.html.kiw`) and VS Code extension support.
- [`node`](modules/node.md) (`github.com/krewire/ui/node`): Streaming `Node` render tree, `SafeURL`, named `Slot`s, layout primitives, and `Kiw()` adapter.
- [`widget`](modules/widget.md) (`github.com/krewire/ui/widget`): Fluent HTML `Widget` builders (`Button`, `Badge`, `Card`, `Alert`, `Modal`, `Table`, `Grid`, `Stack`).
- [`form`](modules/form.md) (`github.com/krewire/ui/form`): Programmatic `Form` builder, `http.Request` binding, validation, and accessible field rendering.
- [`panel`](modules/panel.md) (`github.com/krewire/ui/panel`): Administrative `Panel`, KPI `Stat`, and `Dashboard` widgets.
- [`styles`](modules/styles.md) (`github.com/krewire/ui/styles`): Embedded canonical Pop-Brutalist stylesheet (`styles.CSS`).

## Core Features

- **Pure Go & .kiw Templates**: 100% CGO-free, standalone runtime with Tailwind CSS.
- **Model C Package Distribution**: Distributed as importable Go packages (`github.com/krewire/ui/components`, `github.com/krewire/ui/styles`) with an optional vendor/eject escape hatch.
- **Accessible & Semantic**: Native HTML primitives, ARIA attributes, keyboard navigation, and high-contrast color tokens.
- **Acyclic Leaf Subsystems**: `node`, `widget`, `form`, `styles`, and `kiw` are isolated leaf packages verified by `architecture_test.go`.
