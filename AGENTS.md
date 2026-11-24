# AGENTS.md — Krewire UI

Operational guide for autonomous AI agents working inside the `krewire/ui` repository.

## Architecture

- Pure multi-module monorepo comprising 10 modules (`components`, `core`, `dsl`, `editors`, `form`, `kiw`, `node`, `panel`, `styles`, `widget`).
- No `go.mod` at repository root; coordinated via `go.work`.
- 100% pure-Go standard library only (`CGO_ENABLED=0`).
- All code, comments, docstrings, and commit messages must be strictly in English.

## Quality Gates

Before finishing work:
```bash
gofmt -l .
go vet ./...
go test ./...
CGO_ENABLED=0 go build ./...
```
