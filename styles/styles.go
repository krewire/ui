package styles

import (
	_ "embed"
)

// CSS contains the canonical Krewire unified stylesheet bundle with design tokens,
// CSS variables, base resets, atomic controls, and responsive layout primitives.
//
//go:embed krewire.css
var CSS string

// Compile processes a directory and generates a tree-shaken, purged KrewireCSS stylesheet.
func Compile(root string, opts ...Option) (string, error) {
	c := NewCompiler(opts...)
	return c.CompileDir(root)
}

// CompileClasses generates a tree-shaken KrewireCSS stylesheet for explicit classes.
func CompileClasses(classes []string, opts ...Option) string {
	c := NewCompiler(opts...)
	return c.GenerateCSS(classes)
}
