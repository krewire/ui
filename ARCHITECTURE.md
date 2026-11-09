# Krewire UI (kui) Architecture and Integration Contract

This document defines what `github.com/krewire/ui` (kui) is responsible for, what
it deliberately does not do, and the contract between kui and the `.kiw`
toolchain that renders its components. It is the reference for products (Panel,
Dash, the landing sites) and for the core toolchain in `github.com/krewire/krewire`.

Statements about the core toolchain describe `github.com/krewire/krewire` as of
commit `3bfd916` (2026-11-09). Where the toolchain is the source of a behaviour,
the file is named so it can be re-checked.

---

## 1. Responsibilities

kui owns the **shared presentation foundation**:

| Area | What kui provides | Where |
|---|---|---|
| Components | 66 reusable `.kiw` single-file components (markup + scoped `<style>`, optionally `<script>`) | `components/*.kiw` |
| Component catalog | Name, category, description, props, defaults, dependencies, preview of every component, and starter packs | `components/components.json` (schema: `components.schema.json`), `components.Registry`, `components.Resolve()` |
| Distribution | The components as a standard `embed.FS`, and raw source by name | `components.FS`, `components.Get()` |
| Design tokens & base styles | CSS custom properties (colour, surface, typography, spacing, radius, pop-border/shadow), resets, dark mode, layout primitives | `styles/krewire-ui.css`, `styles.CSS` |
| Layout primitives | `Box`, `Container`, `Flex`, `Grid`, `HStack`, `Section`, `Stack`, `VStack` | `components/` |
| Accessibility baseline | Semantic elements, focus rings, ARIA attributes, reduced motion | components and stylesheet; see `DESIGN_SYSTEM.md` |
| Design specification | Visual language and token values | `DESIGN_SYSTEM.md` |

## 2. Non-goals

kui does **not**:

- **Parse or render `.kiw`.** Parsing and rendering belong to the core toolchain
  (`pkg/kiw`, `pkg/web/ssg`). kui ships template source; it never executes it.
- **Orchestrate builds or ship a CLI.** The standalone `kui` binary was removed
  (commit `a5fb95c`); building, serving and deploying are `kiw`'s job.
- **Depend on any other module.** `go.mod` has no `require` lines and the Go
  packages import only the standard library. This is enforced by
  `TestModuleHasNoExternalDependencies` in `components/contract_test.go`.
- **Contain product business logic.** No data fetching, API calls, auth,
  tenancy, billing, permissions or workflow state. Product screens, navigation
  policies and workflows live in the product repositories (see §6).
- **Provide a typed Go component API.** The programmatic Go UI layer (typed
  renderers, themes, `App` builder) is `github.com/krewire/krewire/pkg/ui`.
  kui is the `.kiw` presentation catalog that layer delegates to
  (`pkg/ui/components/components.go`: `var FS = uicomponents.FS`).
- **Bundle JavaScript frameworks or a Node toolchain.** Interactive behaviour
  prefers native elements (`<details>`, `<dialog>`); any script is small
  vanilla JS.

## 3. Layering

```
 products (panel, dash, krewire.com, krewire.github.io, cloud landing)
   │  .kiw pages/layouts/components            │  optional Go import
   ▼                                           ▼
 core toolchain (github.com/krewire/krewire)   github.com/krewire/ui  (kui)
   kiw CLI ─► pkg/web/ssg ─► pkg/kiw            components  styles
                   │                              ▲
                   └── pkg/ui ── pkg/ui/components┘   (embed.FS delegation)
```

Dependencies point one way: the toolchain imports kui to obtain the built-in
component catalog; **kui never imports the toolchain**. A product that only
needs the catalog and stylesheet imports kui and nothing from
`github.com/krewire/krewire`.

How each sibling repository consumes kui today:

| Consumer | How |
|---|---|
| `krewire` | `pkg/ui/components` re-exports `components.FS`; `pkg/web/ssg/load.go` loads every embedded `.kiw` as a built-in component. |
| `dash`, `cloud` (landing), `krewire.com`, `krewire.github.io` | Use kui components as `.kiw` tags (`<Navbar>`, `<Footer>`, `<VStack>`, `<Button>`, …) rendered by `kiw`; no direct Go import. |
| `panel` | Does not consume kui yet. |

## 4. Go API surface

The public Go API is small and is the part covered by the compatibility policy
in §7:

```go
package components
var FS embed.FS                         // every *.kiw component plus components.json, at the FS root
const ManifestFile = "components.json"
type ComponentMeta struct {             // one components.json entry
	Name, Category, Description, File, Preview string
	Props, Dependencies []string        // derived from the markup
	Defaults map[string]string          // derived: class-helper fallbacks
	Styles, Script bool                 // derived
}
type Pack struct{ Name, Description string; Components []string }
type Manifest struct{ Schema string; Version int; Packs []Pack; Components []ComponentMeta }
var Registry map[string]ComponentMeta   // key == ComponentMeta.Name == file name without .kiw
func List() []string                    // sorted component names
func Get(name string) ([]byte, error)   // "Button" or "Button.kiw"; error wraps fs.ErrNotExist
func LoadManifest() (*Manifest, error)  // decode components.json
func Packs() []Pack                     // starter packs, sorted by name
func Lookup(name string) (ComponentMeta, bool)       // case-insensitive
func Resolve(refs ...string) ([]string, error)       // names and "pack:<name>" plus transitive dependencies

package styles
var CSS string                          // contents of styles/krewire-ui.css
```

`components.json` is the single source of truth for the catalog; its format is
`components/components.schema.json` (JSON Schema Draft-07). Name, category,
description, preview and packs are authored; the other fields are derived from
the markup by `go generate ./components`.

Invariants (tested in `components/components_test.go`,
`components/contract_test.go` and `components/manifest_test.go`):

- Every embedded `.kiw` file has a `Registry` entry and every `Registry` entry
  has an embedded file.
- `Registry` categories are one of: Layout, Atoms, Forms, Navigation, Feedback,
  Overlays, Data Display, Document (the schema enum).
- The derived fields match the markup (`TestManifestInSync`), every dependency
  and pack member is a registered component, and the `minimal`, `base` and
  `dashboard` packs exist.

## 5. The `.kiw` integration contract

### 5.1 Component file

- One component per file; the **file name is the component name**
  (`Button.kiw` → `<Button>`), PascalCase.
- A file contains template markup plus optional `<style>` blocks and optional
  `<script>` blocks. kui components do not use frontmatter.
- A component may nest other components only by their kui names. Every nested
  tag must be a registered kui component (`TestNestedComponentsAreRegistered`).

### 5.2 Template language used by kui components

kui components are written in the `.kiw` template syntax that `pkg/kiw`
desugars to Go `html/template` (`pkg/kiw/kiw.go`, `pkg/kiw/desugar.go`):

| Construct | Meaning |
|---|---|
| `{prop}` | Interpolate a prop (`{{.Prop}}`) |
| `{#if expr}…{:else if expr}…{:else}…{/if}` | Conditionals; kui uses `==`, `!=`, `&&`, `\|\|`, `!` and parentheses in `expr` |
| `{#each items as item}…{/each}` | Loops |
| `<Name attr="…" />`, `<Name …>…</Name>` | Nested component call |
| `{{ … }}` | Plain Go template actions, used where the sugar is not enough |

Go template actions in kui markup may only use text/template keywords and
builtins (`if`, `else`, `end`, `eq`, `printf`, …) plus the renderer helpers
listed in §5.4. This is enforced by `TestComponentsUseOnlyContractHelpers`.

### 5.3 Props, body and slots

- **Props** are the attributes on the component tag. The renderer passes them in
  a map and registers each key with both an upper- and lower-case first letter
  (`pkg/web/ssg/render.go`, `dict`/`componentDict`), so `Variant="primary"` and
  `variant="primary"` are equivalent and a component may read `{variant}` or
  `.Variant`.
- **Value types.** `"true"`, `"false"` and numeric-looking strings (`"1"`,
  `"12"`, `"+2"`) arrive as booleans/numbers, not strings; `attr={expr}` passes
  an expression. Components must therefore not compare a possibly-numeric prop
  against a string literal with `==`; normalise first, as `Heading.kiw` does
  with `eq (printf "%v" .Level) "1"`.
- **Body.** Child content of a paired tag is passed as the `Body` prop and read
  with `{body}`. It is rendered before being passed, so it may contain other
  components.
- **Named slots.** `<div slot="name">…</div>` (also `<section>` or `<template>`)
  and `<slot name="name">…</slot>` inside a component tag become the prop
  `Name`, read with `{name}`. The slot extractor in `pkg/kiw` ends a
  `div`/`section`/`template` slot at the first closing `div`/`section`/`template`
  tag, so content inside such a slot must not itself contain those elements; use
  `<slot name="…">` when it does.
- A component's props and slots are documented by its markup (`{#if prop}` /
  `{prop}`) and by the per-component pages on krewire.com.

### 5.4 Renderer requirements

A renderer that wants to render kui components must provide:

1. The `.kiw` template sugar in §5.2 and the prop/body/slot rules in §5.3.
2. These template functions:
   - `class` — class composition: `{{class "btn" .Variant "primary" .Size "md" .Class}}`
     emits the base class, then `base-<value or default>` for each
     (value, default) pair, then trailing raw tokens
     (reference implementation: `ClassHelper` in `pkg/web/ssg/class.go`).
   - `assetLinks` — emits the site's stylesheet/script tags; used only by
     `Head.kiw`.
   - `component` and `dict` — component invocation; the desugared form of
     every nested tag.
3. Override precedence: a project's own `components/<Name>.kiw` replaces the kui
   component of the same name (`pkg/web/ssg/load.go` loads kui first, then the
   project directory).
4. Style handling: each component's `<style>` is emitted only when the component
   is used, scoped under `[data-kiw-component="Name"]`; rules on `:root`,
   `@keyframes` and `@font-face` stay global (`pkg/web/ssg/css.go`). The
   renderer adds `data-kiw-component="Name"` to the component's first element
   when the markup does not already carry it (`pkg/web/ssg/scope.go`).

The minimum toolchain packages for rendering are therefore
`github.com/krewire/krewire/pkg/kiw` (parser) and
`github.com/krewire/krewire/pkg/web/ssg` (renderer). A product that does not
render `.kiw` itself (for example, it only serves `styles.CSS` or ships the
component sources) needs neither.

**Known gap.** `pkg/kiw` extracts inline `<script>` blocks out of the markup,
and the `pkg/web/ssg` component loader keeps only the markup and styles. The
inline scripts in `CodeWindow.kiw` and `Terminal.kiw` are therefore not emitted
by that loader; those components must stay usable without them.

## 6. Shared foundation vs product-specific code

Rule: kui contains presentation only. Anything that knows about a specific
product's data, routes, permissions or workflows belongs to that product.

Audit of the current catalog (no component fetches data, calls APIs, stores
state in `localStorage`, or opens sockets). The following components carry
**Krewire brand content as defaults**:

| Component | Krewire-specific default | How a product avoids it |
|---|---|---|
| `Brand` | Name "Krewire" and the Krewire logo SVG | `name`, `icon`, or a body |
| `Navbar` | Ecosystem product links in the mobile drawer when `showEcosystemMenu` is set | body / `drawer` slot, `navLinks`; leave `showEcosystemMenu` unset |
| `NavMenu` | Krewire product/workload dropdowns (`showEcosystemMenu`, `showProductsDropdown`) | do not set those flags; pass a body |
| `Footer` | The `contributors` variant and the default variant (any other value, e.g. `full`) without `columns` list Krewire products and repositories; copyright defaults to "Krewire Contributors" | `variant="simple"` with `brandName`, `columns`, or a body |
| `Ecosystem` | Hard-coded list of Krewire repositories | `header` slot and a body |
| `SideMenu` | Default table of contents of the Krewire docs | pass a body |
| `CodeWindow` | Sample code importing `github.com/krewire/krewire/pkg/web` | pass your own code |

These are brand defaults, not business logic, and every one can be overridden
through props or slots. They stay in kui for now because Krewire's own sites use
them (`<Navbar>` and `<Footer>` in the dash, cloud, krewire.com and
krewire.github.io layouts, mostly with the default/`full` footer; `<NavMenu>` in
krewire.com). Policy going forward:

- No new product- or brand-specific content is added to kui components.
- Product-specific UI is written as a product component that composes kui
  primitives (see `examples/productapp/site/components/ServerCard.kiw`), or as a
  product override of a kui component (§5.4, item 3).
- Moving the brand defaults out of kui is a breaking change for the sites above
  and follows §7.

## 7. Compatibility expectations

**Current state.** The module path is `github.com/krewire/ui`; `go.mod`
declares `go 1.27.1`; the repository has **no release tags yet**, so consumers
resolve it by pseudo-version. Consumers in the ecosystem (`krewire`, `panel`,
`dash`) also declare `go 1.27.1`.

**Versioning.** Releases follow Semantic Versioning with Go module tags
(`vMAJOR.MINOR.PATCH`). Until `v1.0.0`, minor releases may contain breaking
changes, which are called out in the release notes and commit messages (`!` or
a `BREAKING CHANGE:` footer). From `v1.0.0`, breaking changes require a new
major version.

**What counts as the public contract.**

| Surface | Breaking change examples |
|---|---|
| Go API (§4) | Removing or changing `FS`, `Registry`, `ComponentMeta`, `List`, `Get`, `CSS` |
| Component names | Renaming or removing a `.kiw` file |
| Props and slots | Removing or renaming a prop/slot, changing its default, changing a variant name |
| Rendered hooks | Removing or renaming a CSS class or `data-*` attribute that products target |
| Design tokens | Removing or renaming a CSS custom property in `styles/krewire-ui.css` |
| Renderer requirements (§5.4) | Requiring a new template helper or DSL construct |

Adding components, props, slots, variants or tokens is non-breaking. Visual
refinements that keep names and structure are non-breaking.

**Go version.** The `go` directive moves only together with the core toolchain
and is never raised in a patch release.

**Copy-and-own consumers.** Projects that copy `.kiw` files into their own
`components/` directory own those copies; upstream changes reach them only when
they copy again. `kiw add ui/<Name>` records each copy in
`components/.kui-lock.json` (kui version and SHA-256), refreshes only
unmodified copies, and `kiw ui diff <Name>` compares a copy with upstream; see
"Updating Customized Components" in README.md. Projects that rely on the built-in catalog through `kiw` pick
up changes whenever their `github.com/krewire/ui` version changes.

## 8. Composition patterns and examples

| Pattern | Example |
|---|---|
| Layout composition with primitives | `examples/productapp/site/pages/index.kiw` (`Container` → `VStack` → `HStack` → `Stat`) |
| Body and named slots | `examples/productapp/site/layouts/Base.kiw` (`Navbar` body + `actions` slot) and `ServerCard.kiw` (`Card` body + `footer` slot) |
| Product component wrapping kui | `examples/productapp/site/components/ServerCard.kiw` |
| Product override of kui | `Example_productOverride` in `components/example_test.go`; `Catalog` in `examples/productapp/main.go` |
| Theming by tokens | Override custom properties (`--primary`, `--radius`, …) on `:root` or a container; see `DESIGN_SYSTEM.md` |
| Go consumption without the toolchain | `examples/productapp` (imports only kui and the standard library), `ExampleCSS` in `styles/example_test.go` |

All examples run under `go test ./...`. The example site in
`examples/productapp/site` is a regular `.kiw` site; it was checked to render
with `pkg/web/ssg` at the toolchain commit named at the top of this document.
