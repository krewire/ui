# Contributing to Krewire UI

Guidelines for contributing components, styles, and primitives to **Krewire UI** (`github.com/krewire/ui`).

> All code, comments, docstrings, commit messages, and documentation are strictly in **English**.

## Invariants

1. **Pure Go & Zero CGO**: All modules must compile with `CGO_ENABLED=0`.
2. **Modular Architecture**: Code belongs to its respective module (`core/`, `components/`, `styles/`, etc.).
3. **Accessibility**: Components must support native HTML semantics, ARIA landmarks, and high contrast.

## Quality Gates

```bash
gofmt -l .
go vet ./...
go test ./...
CGO_ENABLED=0 go build ./...
```
