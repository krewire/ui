package styles

import (
	_ "embed"
)

// CSS contains the canonical Krewire unified stylesheet bundle with design tokens,
// CSS variables, base resets, atomic controls, and responsive layout primitives.
//
//go:embed krewire.css
var CSS string
