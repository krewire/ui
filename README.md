# Krewire UI (kui)

Composable, zero-dependency `.kiw` design system, layout & UI components for the Krewire ecosystem.

Inspired by the **shadcn/ui** philosophy: you do not install an opaque `node_modules` directory or inherit a rigid black-box framework. Instead, you **own your components**. You copy `.kiw` component files directly into your project's `components/` directory, adapting the markup, logic, and styling to fit your exact product requirements.

[![Go Reference](https://pkg.go.dev/badge/github.com/krewire/ui.svg)](https://pkg.go.dev/github.com/krewire/ui)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Why Krewire UI?

- **Zero Runtime Bloat**: 100% Pure Go and `.kiw` component templates with zero npm, node, or webpack dependencies.
- **Copy & Own (shadcn Model)**: Components live in your project's `./components` directory. No upstream lock-in or breaking library upgrades.
- **Signature Pop-Brutalist Aesthetic**: Built-in neo-brutalist pop shadows, crisp outlines, vibrant accents, and high-contrast surface tokens.
- **Unified Design Tokens**: Comprehensive typography scale, 4px/8px grid spacing, surface hierarchy, and dark mode tokens.
- **WCAG Accessibility**: Full support for native HTML primitives (`<details>`, `<dialog>`, semantic tags), ARIA landmarks, focus rings, and reduced motion.
- **Go Embed Native**: Import `github.com/krewire/ui/components` directly in Go binaries via standard library `embed.FS`.

---

## Getting Started

### 1. Add Components to Your Project

Install components with the `kiw` CLI. It copies the `.kiw` source that ships
with your `kiw` build into `./components/`, together with any kui components it
nests (for example `Navbar` brings `Brand`, `Hamburger`, `NavMenu`, …):

```bash
kiw ui list                    # browse the catalog, packs, and what you already have
kiw ui info Card               # props, defaults, dependencies and a preview snippet
kiw add ui/button ui/card      # one or more components (names are case-insensitive)
kiw add ui:minimal             # a starter pack
```

Starter packs:

| Pack | Components |
|---|---|
| `minimal` | Button, Card, Input |
| `base` | Container, Button, Card, Input |
| `dashboard` | Table, Stat, Sidebar, Modal |

Without `kiw`, copy the files by hand from a checkout of this repository
(`cp path/to/ui/components/Button.kiw ./components/`).

### 2. Include the Design System Stylesheet

Include the canonical stylesheet containing design tokens, base resets, and layout primitives in your page layout:

```html
<link rel="stylesheet" href="/styles/krewire-ui.css">
```

### 3. Use in `.kiw` Pages

Child content becomes the component body (`{body}`); `<slot name="…">` or
`<div slot="…">` passes a named slot. See
[ARCHITECTURE.md](ARCHITECTURE.md#53-props-body-and-slots) for the full contract.

```html
<Navbar brandName="My App">
  <a href="/docs" class="nav-link">Documentation</a>
  <a href="/pricing" class="nav-link">Pricing</a>
  <slot name="actions">
    <Button variant="primary" size="sm">Get Started</Button>
  </slot>
</Navbar>

<Container>
  <VStack gap="lg" align="center">
    <Heading level="1">Modern Go Web Applications</Heading>
    <Text size="lg" color="muted">Zero-dependency, composable UI for Go.</Text>

    <HStack gap="md">
      <Button variant="primary">Deploy Now</Button>
      <Button variant="ghost">Learn More</Button>
    </HStack>

    <Stat label="Total Deployments" value="1,429" change="+18.4%" trend="up" subtitle="Past 30 days" />
  </VStack>
</Container>
```

---

## Updating Customized Components

Copied components are yours: edit them freely. `kiw` tracks what it copied so it
can tell your edits apart from upstream changes:

- Each install records the kui version and the SHA-256 of the copied file in
  `components/.kui-lock.json`. Commit this file with your components.
- `kiw ui list` shows the state of every installed component: `up to date`,
  `modified` (you changed it), `update available` (kui changed it), or both.
- `kiw ui diff <Name>` shows a unified diff from the upstream version to your
  local copy (uses `git diff --no-index`).
- `kiw add ui/<Name>` again updates a component **only if you have not modified
  it**. A modified component is never overwritten silently; the command stops and
  points you to `kiw ui diff`.

Recommended update workflow for a customized component:

1. `kiw ui diff Card` — review what you changed relative to upstream.
2. Keep your changes small and local: prefer props, slots and CSS custom
   properties over editing markup, so upstream updates stay easy to merge.
3. To take the new upstream version, commit your copy, run
   `kiw add ui/card --force`, then re-apply your edits with
   `git diff` / `git checkout -p` against the commit you just made.
4. To keep your version, do nothing: the lock entry keeps reporting
   `modified` so the decision stays visible.

Upgrading `kiw` (or the `github.com/krewire/ui` version a project builds with)
is what brings new upstream component versions.

---

## Component Manifest

`components/components.json` catalogs every component: name, category,
description, `file`, `props` (lower-camel; `body` is the child content),
class-helper `defaults`, nested-component `dependencies`, whether it ships
`styles`/`script`, and a `preview` snippet. It also defines the starter
`packs`. The format is described by
[`components/components.schema.json`](components/components.schema.json)
(JSON Schema Draft-07).

The derived fields are computed from the markup; after editing a component run:

```bash
go generate ./components     # refresh components.json
go test ./...                # TestManifestInSync fails if it is stale
```

---

## Go Module / Embed Usage

For Go backend services or applications that want to embed components directly:

```go
package main

import (
	"fmt"
	"github.com/krewire/ui/components"
	"github.com/krewire/ui/styles"
)

func main() {
	// Access embedded filesystem containing all 66 .kiw components
	fs := components.FS

	// Read a specific component
	buttonKiw, err := components.Get("Button")
	if err != nil {
		panic(err)
	}

	// Catalog metadata and starter packs from components.json
	meta := components.Registry["Button"] // props, defaults, dependencies, preview
	deps, _ := components.Resolve("pack:minimal", "Navbar")
	fmt.Println(meta.Preview, deps)

	// Access canonical CSS stylesheet
	css := styles.CSS
	fmt.Printf("Loaded Button (%d bytes), CSS (%d bytes)\n", len(buttonKiw), len(css))
}
```

kui depends only on the Go standard library; importing it does not pull in the
`kiw` CLI or `github.com/krewire/krewire`.

---

## Architecture, Contract & Examples

- [ARCHITECTURE.md](ARCHITECTURE.md): responsibilities and non-goals, the `.kiw`
  integration contract (props, body, slots, renderer requirements), the
  boundary with product code, and compatibility/versioning expectations.
- [DESIGN_SYSTEM.md](DESIGN_SYSTEM.md): visual language and token values.
- Runnable Go examples: `components/example_test.go`, `styles/example_test.go`
  (`go test ./...` checks their output).
- [`examples/productapp`](examples/productapp): a product application that
  consumes kui on its own, with a `.kiw` site composing kui primitives and a
  product-owned `ServerCard` component.

---

## Design System Tokens

Krewire UI defines a standardized set of CSS custom properties in `styles/krewire-ui.css`:

### Color Palette & Surfaces
- Primary: `--primary` (`#39D353`), `--primary-content` (`#0B1F3B`), `--primary-hover`, `--primary-subtle`
- Secondary: `--secondary` (`#00D1C1` in light, `#112646` in dark), `--secondary-content`
- Accent: `--accent` (`#FF3B2E`), `--accent-content` (`#ffffff`)
- Status / Feedback: `--success` (`#39D353`), `--info` (`#00D1C1`), `--warning` (`#ffb703`), `--error` (`#FF3B2E`)
- Surfaces: `--bg` (`#FAF8F4` in light, `#0B1F3B` in dark), `--bg-2`, `--bg-3`, `--bg-subtle`
- Typography: `--fg` (`#0B1F3B` in light, `#FAF8F4` in dark), `--fg-muted`, `--fg-subtle`
- Borders: `--border`, `--border-subtle`, `--border-focus`

### Neo-Brutalist Pop Tokens
- Border: `--pop-border: 2.5px solid var(--border);`
- Thin Border: `--pop-border-thin: 2px solid var(--border);`
- Standard Shadow: `--pop-shadow: 4px 4px 0px var(--border);`
- Small Shadow: `--pop-shadow-sm: 2px 2px 0px var(--border);`
- Large Shadow: `--pop-shadow-lg: 6px 6px 0px var(--border);`

### Typography Scale
- Fonts: `--ui-font-sans`, `--ui-font-display`, `--ui-font-mono`
- Type Sizes: `--ui-text-2xs` (10px), `--ui-text-xs` (11px), `--ui-text-sm` (13px), `--ui-text-base` (14.5px), `--ui-text-lg` (17px), `--ui-text-xl` (20px), `--ui-text-2xl` (24px), `--ui-text-3xl` (30px), `--ui-text-4xl` (36px), `--ui-text-5xl` (48px)
- Line Heights: `--ui-leading-tight` (1.25), `--ui-leading-snug` (1.375), `--ui-leading-normal` (1.5), `--ui-leading-relaxed` (1.625)

### Spacing Scale (4px/8px Grid)
- `--ui-space-1` (4px), `--ui-space-2` (8px), `--ui-space-3` (12px), `--ui-space-4` (16px), `--ui-space-6` (24px), `--ui-space-8` (32px), `--ui-space-12` (48px), `--ui-space-16` (64px)

---

## Component Catalog (66 Components)

### 1. Layout Primitives (8)
- `Box`: Generic layout container with customizable padding, margin, and borders.
- `Container`: Responsive max-width page container.
- `Flex`: Flexible box layout container with direction, wrap, and alignment.
- `Grid`: Responsive CSS grid container with customizable columns and gaps.
- `HStack`: Horizontal flex stack with configurable spacing and alignment.
- `Section`: Semantic section wrapper with vertical rhythm.
- `Stack`: Directional stack wrapper with uniform spacing.
- `VStack`: Vertical flex stack with configurable spacing and alignment.

### 2. Atoms & Primitives (12)
- `Avatar`: Image or initials avatar with size variants and status indicators.
- `Badge`: Compact label for status, counts, or categorization.
- `Brand`: App brand logo and name display.
- `Button`: Interactive button with primary, secondary, ghost, outline, and error variants.
- `Code`: Inline code badge with monospace typography.
- `Divider`: Horizontal or vertical separator line.
- `Heading`: Semantic heading element (h1-h6) with consistent typography tokens.
- `Kbd`: Keyboard shortcut keycap indicator.
- `Link`: Accessible hyperlink with hover states and external indicator.
- `Progress`: Determinate or indeterminate progress bar indicator.
- `Spinner`: Loading spinner indicator with speed and color variants.
- `Text`: Versatile typography text atom with muted, lead, and size props.

### 3. Forms & Inputs (10)
- `ButtonGroup`: Segmented button group for view switching or filter toggles.
- `Checkbox`: Accessible checkbox input with label and description support.
- `Form`: Form layout wrapper with accessible grouping and validation slots.
- `Input`: Text input with prefix/suffix icons, helper text, and error states.
- `Radio`: Radio selection control for mutually exclusive options.
- `SearchInput`: Search input box with shortcut badge and clear button.
- `Select`: Dropdown selection menu with custom styled options.
- `Slider`: Range slider input with thumb, track, and value indicator.
- `Textarea`: Multi-line text input with autoresize and character counter.
- `Toggle`: Binary switch toggle for boolean settings.

### 4. Navigation (11)
- `Breadcrumb`: Hierarchical page breadcrumb navigation trail.
- `BreadcrumbItem`: Individual crumb segment within a Breadcrumb navigation.
- `Hamburger`: Animated mobile navigation toggle button.
- `LangSwitch`: Multi-language locale switcher dropdown.
- `Navbar`: Top navigation bar with brand, menu links, slots, and responsive drawer.
- `NavMenu`: Composable navigation item list with slot for active page tracking.
- `Pagination`: Responsive pagination bar with page numbers, prev/next, and ellipsis.
- `Sidebar`: Full-height vertical collapsible side navigation container.
- `SideMenu`: Structured sidebar links with category groupings and badges.
- `Tabs`: Tabbed content switcher with underline or pill styling.
- `ThemeSwitch`: Light and dark mode theme toggle switch.

### 5. Feedback & Alerts (4)
- `Alert`: Contextual alert box with info, success, warning, and error types.
- `Banner`: Top or bottom announcement banner with action link and dismiss button.
- `Callout`: Prominent callout banner for highlights, tips, or quotes.
- `Toast`: Floating notification toast banner with auto-dismiss and status variants.

### 6. Overlays (7)
- `ActionMenu`: Contextual kebab action menu for table rows and card action items.
- `Collapse`: Expandable accordion disclosure block with smooth animation.
- `Drawer`: Slide-out panel overlay anchored to left or right screen edge.
- `Dropdown`: Composable dropdown menu with trigger, item shortcuts, and dividers.
- `Modal`: Accessible dialog modal overlay with header, body, and footer slots.
- `Popover`: Floating rich content popover with header and trigger button.
- `Tooltip`: Contextual hover tooltip bubble.

### 7. Data Display & Complex (13)
- `Accordion`: Multi-item collapsible disclosure accordion with smooth chevron rotation.
- `Card`: Content card container with header, media, body, and footer slots.
- `CodeWindow`: Syntax-highlighted code display window with copy button and titlebar.
- `Ecosystem`: Interactive ecosystem showcase grid with category filtering.
- `Empty`: Zero-data state illustration with title, message, and call-to-action button.
- `Footer`: Full-width site footer with customizable link columns and copyright.
- `List`: Styled bulleted or numbered list container with custom markers.
- `Skeleton`: Content placeholder loading skeleton with subtle shimmer animation.
- `Stat`: Metric KPI display card with label, value, trend badge, and subtitle.
- `Table`: Responsive data table with striped, bordered, and hover states.
- `Tag`: Interactive tag chip with optional removable button and filter state.
- `Terminal`: Retro-modern CLI terminal window display with animated prompts.
- `Timeline`: Vertical event history and changelog milestone timeline.

### 8. Document (1)
- `Head`: HTML head document tags injector for titles, meta, and styles.

---

## License

MIT License — free for personal, commercial, and enterprise use.
