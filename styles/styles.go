package styles

import (
	_ "embed"
)

// CSS contains the canonical Krewire UI stylesheet with design tokens,
// CSS variables, base resets, and responsive layout primitives.
//
//go:embed krewire-ui.css
var CSS string
