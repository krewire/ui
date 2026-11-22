# Contributing to Krewire UI

Guidelines for contributing to the Krewire Pop-Brutalist UI framework and `.kiw` component library.

## Architectural Boundaries

- **Leaf Packages**: `node`, `widget`, `form`, `styles`, and `kiw` must never import the root `ui` facade.
- **Pure Go**: Always compile with `CGO_ENABLED=0`. Zero Node.js or npm dependencies in core or component packages.
- **Model C Invariant**: Components are distributed as importable Go packages. Never require manual copying of files for basic framework usage.

## Component Design Standards

1. Use semantic HTML tags (`<button>`, `<dialog>`, `<details>`, `<nav>`).
2. Provide ARIA labels and focus rings for accessibility.
3. Support light and dark mode tokens seamlessly.
4. Keep component templates self-contained and declarative.

## Verification Checklist

Before opening a pull request or committing changes:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
