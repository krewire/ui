// Package ui provides the primary public facade, application builder, layout
// engine, and theme primitives for the Krewire Pop-Brutalist UI framework.
//
// # Architectural Model & Subsystem Boundaries
//
// The Krewire UI framework is structured as an acyclic, modular hierarchy
// following the 6-Layer Pyramid Architecture (ADR-E8P9D) and Model C package
// distribution (ADR-K7P9L). The root package serves as the unified, ergonomic
// entrypoint for application assembly while domain capabilities are cleanly
// isolated in leaf packages:
//
//   - ui: Primary public facade, AppBuilder, Layout, Theme, Palette, and Registry.
//   - ui/components: Shipped .kiw component assets, embedded filesystem, and catalog.
//   - ui/form: Accessible form model, typed inputs, CSRF bindings, and validation.
//   - ui/node: High-performance virtual DOM / render tree and semantic HTML nodes (leaf).
//   - ui/widget: Standalone HTML widgets (Button, Badge, Card, Alert, Table, Modal).
//   - ui/panel: Dashboard layout composition, metric cards, and collapsible panels.
//   - ui/styles: Embedded Pop-Brutalist CSS tokens and theme variables (leaf).
//   - ui/kiw: Go-native component compiler, template desugaring, and diagnostics.
//   - ui/dsl: Domain-specific language contracts and template bindings.
//   - ui/editors: Language server and editor integration protocols.
//
// # Invariants
//
// 1. Leaf Subsystem Independence: Leaf packages (node, form, widget, styles, kiw)
// must never import the root ui facade or higher-level consumers.
//
// 2. Zero npm / Node.js Runtime: All rendering, compilation, and styling is
// 100% pure Go (CGO_ENABLED=0).
//
// 3. Ergonomic Facade Stability: The root package provides stable, backward-compatible
// composition interfaces without breaking downstream consumers across the ecosystem.
package ui
