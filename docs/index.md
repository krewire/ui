# Krewire UI — Developer & Contributor Documentation

Welcome to the internal engineering documentation for **`github.com/krewire/ui`**, the pop-brutalist `.kiw` design system, component catalog, and layout engine for the Krewire digital SDLC ecosystem.

## Documentation Index

- [Architecture & Modular Subsystems](architecture.md): 6-Layer pyramid integration, Model C package distribution, subsystem boundaries (`node`, `widget`, `form`, `styles`, `panel`, `kiw`), and layering invariants.
- [Design System & Tokens](design-system.md): Pop-Brutalist design tokens, typography scales, spacing grids, high-contrast surfaces, and dark mode palettes.
- [Getting Started](getting-started.md): Installation, package-based imports, vendoring components (`krew ui vendor`), layouts, and page integration.
- [Contributing Guide](contributing.md): Quality gates, adding components, tests, and contribution standards.

## Core Features

- **Pure Go & .kiw Templates**: 100% CGO-free, zero npm/Node.js runtime bloat.
- **Model C Package Distribution**: Distributed as importable Go packages (`github.com/krewire/ui/components`, `github.com/krewire/ui/styles`) with an optional vendor/eject escape hatch.
- **Accessible & Semantic**: Native HTML primitives, ARIA attributes, keyboard navigation, and high-contrast color tokens.
- **Acyclic Leaf Subsystems**: `node`, `widget`, `form`, `styles`, and `kiw` are isolated leaf packages verified by `architecture_test.go`.
