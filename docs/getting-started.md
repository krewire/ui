# Getting Started with Krewire UI

This guide explains how to install, import, and use components from `github.com/krewire/ui`.

## 1. Import as a Go Module (Model C)

In your project's `go.mod`, import `github.com/krewire/ui`:

```bash
krew add github.com/krewire/ui
```

In `.kiw` templates, import components directly:

```html
<script>
  import { Button, Card, Navbar } from '@krewire/ui';
</script>

<Card>
  <h2>Welcome to Krewire</h2>
  <p>Pop-brutalist UI components in pure Go.</p>
  <Button variant="primary">Get Started</Button>
</Card>
```

## 2. Optional Component Vendoring (Eject Escape Hatch)

To customize or own component code locally:

```bash
krew ui vendor Card
# or: krew ui eject Card
```

This copies `Card.kiw` into your local `components/` directory while tracking upstream versions in `.kui-lock.json`.

## 3. Include Design System Stylesheet

Include the clean, modern CSS tokens in your HTML head or layout:

```html
<link rel="stylesheet" href="/styles/krewire.css">
```

Or embed `github.com/krewire/ui/styles` in pure Go binaries via `styles.FS`.

## 4. Running Quality Gates

```bash
cd ui
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
