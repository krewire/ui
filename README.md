# Krewire UI (`github.com/krewire/ui`)

Composable, zero-dependency `.kiw` design system, layout, and UI component catalog for the Krewire digital SDLC ecosystem.

Inspired by the **shadcn/ui** philosophy: components stay transparent and forkable, but distribution is **package-based by default**. You import modular Go packages (`github.com/krewire/ui/core`, `styles`, `components`, etc.) with embedded `.kiw` templates, and only copy a component's source into your project's `components/` directory when you explicitly choose to fork it via `krew ui vendor` / `krew ui eject`.

[![Go Reference](https://pkg.go.dev/badge/github.com/krewire/ui.svg)](https://pkg.go.dev/github.com/krewire/ui)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Modular Packages (Pure Multi-Module Layout)

Krewire UI is organized into **10 independent Go modules**:

- **`core` (`github.com/krewire/ui/core`)**: Theme management, component registry, app builder, and layout tokens.
- **`components` (`github.com/krewire/ui/components`)**: 66 Pop-Brutalist `.kiw` templates embedded via standard library `embed.FS`.
- **`styles` (`github.com/krewire/ui/styles`)**: Scoped CSS tokens and stylesheets (`krewire.css`).
- **`dsl` (`github.com/krewire/ui/dsl`)**: Domain Specific Language parser for `.kiw` templates.
- **`kiw` (`github.com/krewire/ui/kiw`)**: Template compiler and tag lexer.
- **`widget` (`github.com/krewire/ui/widget`)**: Complex interactive UI widgets.
- **`panel` (`github.com/krewire/ui/panel`)**: Dockable panels and dashboard containers.
- **`node` (`github.com/krewire/ui/node`)**: AST DOM node representations.
- **`form` (`github.com/krewire/ui/form`)**: Form controls and validation bindings.
- **`editors` (`github.com/krewire/ui/editors`)**: Code and markdown editor primitives.

---

## Why Krewire UI?

- **Pure Go & Tailwind CSS**: 100% Pure Go and `.kiw` component templates with seamless Tailwind CSS utility styling.
- **Package-Based by Default**: Import modules directly via Go workspaces; no node_modules required for standard Go apps.
- **Explicit Vendor/Eject Escape Hatch**: Fork individual components into `./components` with `krew ui vendor <Name>` only when local custom modifications are needed.
- **Signature Pop-Brutalist Aesthetic**: Neo-brutalist pop shadows, crisp outlines, vibrant accents, and high-contrast surface tokens.
- **Go Embed Native**: Import component templates directly into single static binaries via `embed.FS`.

---

## Quality Gates

Before submitting changes:
```bash
gofmt -l .
go vet ./...
go test ./...
CGO_ENABLED=0 go build ./...
```

## License

[MIT License](LICENSE) © 2026 Krewire Contributors
