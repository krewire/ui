# Krewire UI (`kui`)

Composable, zero-dependency `.kiw` layout & UI components for the Krewire ecosystem.

Inspired by **shadcn/ui**: you do not install an opaque node_modules or rigid black-box library. Instead, you **own your components**. You copy-paste or install the `.kiw` files directly into your project via the `kui` CLI, and adapt the markup and CSS styles as your project demands.

[![Go Reference](https://pkg.go.dev/badge/github.com/krewire/ui.svg)](https://pkg.go.dev/github.com/krewire/ui)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## Why Krewire UI?

- **Zero Runtime Dependencies**: Pure Go and `.kiw` templates with zero external npm/node dependencies.
- **Own Your Code**: Components live directly in your project's `./components` directory. No breaking upstream lock-ins.
- **Native `.kiw` DSL**: Native support for scoped `<style>`, slots, conditional rendering (`{#if}`), loops (`{#each}`), and props.
- **Integrated Design System**: Complete set of CSS custom properties (colors, typography, geometry, shadows) with built-in Light and Dark themes.
- **Dual Distribution**:
  1. **CLI / Copy-Paste Mode**: Add `.kiw` files into your project via `kui add <component>`.
  2. **Go Embedded FS Mode**: Import `github.com/krewire/ui/components` directly in Go apps via `embed.FS`.

---

## Installation & Setup

### 1. Install the `kui` CLI

```bash
go install github.com/krewire/ui/cmd/kui@latest
```

### 2. Initialize in Your Project

In any Krewire project (`site`, `app`, `book`, etc.):

```bash
kui init
```

This generates:
- `kui.json`: Configuration defining component and CSS directories.
- `components/`: Target folder for your `.kiw` components.
- `static/css/krewire-ui.css`: Design system CSS variables and layout primitives.

### 3. Add Components

Add only the components you need:

```bash
# Add specific components
kui add Button Navbar HStack VStack Card

# Or install the full suite (51 components)
kui add --all
```

---

## CLI Command Reference

- `kui init`: Initialize Krewire UI and download the base stylesheet.
- `kui list`: Browse all available components, categories, and descriptions.
- `kui add <names...>`: Copy `.kiw` components into your `./components/` directory.
- `kui info <name>`: View component documentation and inspect template preview.
- `kui version`: Display CLI version.

---

## Using Components in `.kiw`

Once a component is added to your project, you can use it anywhere in your `.kiw` pages or layouts:

```html
<Navbar brand="My App" current="/">
  <slot name="menu">
    <a href="/docs" class="nav-link">Documentation</a>
    <a href="/pricing" class="nav-link">Pricing</a>
  </slot>
  <slot name="actions">
    <Button variant="primary" size="sm">Get Started</Button>
  </slot>
</Navbar>

<Container>
  <VStack gap="lg" align="center">
    <Heading level="1">Welcome to Krewire UI</Heading>
    <Text size="lg" muted="true">Zero-dependency, composable UI for Go.</Text>

    <HStack gap="md">
      <Button variant="primary">Deploy Now</Button>
      <Button variant="ghost">Learn More</Button>
    </HStack>
  </VStack>
</Container>
```

---

## Go Module / Embed Usage

For Go backend services or applications that want to mount components without copying individual files:

```go
package main

import (
	"fmt"
	"github.com/krewire/ui/components"
	"github.com/krewire/ui/styles"
)

func main() {
	// Access embedded filesystem containing all 51 .kiw components
	fs := components.FS

	// Read a specific component
	buttonKiw, err := components.Get("Button")
	if err != nil {
		panic(err)
	}

	// Access canonical CSS stylesheet
	css := styles.CSS
	fmt.Printf("Loaded Button (%d bytes), CSS (%d bytes)\n", len(buttonKiw), len(css))
}
```

---

## Component Catalog (51 Components)

### Layout Primitives
- `Box`: Generic layout container with customizable padding, margin, and borders.
- `Container`: Responsive max-width page container.
- `Flex`: Flexible box layout container with direction, wrap, and alignment.
- `Grid`: Responsive CSS grid container with customizable columns and gaps.
- `HStack`: Horizontal flex stack with configurable spacing and alignment.
- `Section`: Semantic section wrapper with vertical rhythm.
- `Stack`: Directional stack wrapper with uniform spacing.
- `VStack`: Vertical flex stack with configurable spacing and alignment.

### Atoms & Primitives
- `Avatar`: Image or initials avatar with size variants and status indicators.
- `Badge`: Compact label for status, counts, or categorization.
- `Brand`: App brand logo and name display.
- `Button`: Interactive button with primary, secondary, ghost, and danger variants.
- `Code`: Inline code badge with monospace typography.
- `Divider`: Horizontal or vertical separator line.
- `Heading`: Semantic heading element (h1-h6) with consistent typography tokens.
- `Kbd`: Keyboard shortcut keycap indicator.
- `Link`: Accessible hyperlink with hover states and external indicator.
- `Progress`: Determinate or indeterminate progress bar indicator.
- `Spinner`: Loading spinner indicator with speed and color variants.
- `Text`: Versatile typography text atom with muted, lead, and size props.

### Forms & Inputs
- `Checkbox`: Accessible checkbox input with label and description support.
- `Form`: Form layout wrapper with accessible grouping and validation slots.
- `Input`: Text input with prefix/suffix icons, helper text, and error states.
- `Radio`: Radio selection control for mutually exclusive options.
- `Select`: Dropdown selection menu with custom styled options.
- `Textarea`: Multi-line text input with autoresize and character counter.
- `Toggle`: Binary switch toggle for boolean settings.

### Navigation
- `Breadcrumb`: Hierarchical page breadcrumb navigation trail.
- `BreadcrumbItem`: Individual crumb segment within a Breadcrumb navigation.
- `Hamburger`: Animated mobile navigation toggle button.
- `LangSwitch`: Multi-language locale switcher dropdown.
- `Navbar`: Top navigation bar with brand, menu links, slots, and responsive drawer.
- `NavMenu`: Composable navigation item list with slot for active page tracking.
- `Sidebar`: Full-height vertical collapsible side navigation container.
- `SideMenu`: Structured sidebar links with category groupings and badges.
- `Tabs`: Tabbed content switcher with underline or pill styling.
- `ThemeSwitch`: Light and dark mode theme toggle switch.

### Feedback & Status
- `Alert`: Contextual alert box with info, success, warning, and error types.
- `Callout`: Prominent callout banner for highlights, tips, or quotes.

### Overlays
- `Collapse`: Expandable accordion disclosure block with smooth animation.
- `Drawer`: Slide-out panel overlay anchored to left or right screen edge.
- `Modal`: Accessible dialog modal overlay with header, body, and footer slots.
- `Tooltip`: Contextual hover tooltip bubble.

### Data Display & Presentation
- `Card`: Content card container with header, media, body, and footer slots.
- `CodeWindow`: Syntax-highlighted code display window with copy button and titlebar.
- `Ecosystem`: Interactive ecosystem showcase grid with category filtering.
- `Footer`: Full-width site footer with customizable link columns and copyright.
- `List`: Styled bulleted or numbered list container with custom markers.
- `Table`: Responsive data table with striped, bordered, and hover states.
- `Terminal`: Retro-modern CLI terminal window display with animated prompts.

### Document
- `Head`: HTML head document tags injector for titles, meta, and styles.

---

## License

MIT License — free for personal, commercial, and enterprise use.
