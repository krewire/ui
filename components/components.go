// Package components embeds the Krewire UI (kui) .kiw component catalog.
package components

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

// FS embeds all built-in .kiw reusable and atomic UI components and the
// components.json manifest that catalogs them.
//
//go:embed *.kiw components.json
var FS embed.FS

// Registry maps each component name to its metadata. It is decoded from the
// embedded components.json manifest, which is the single source of truth.
var Registry = mustRegistry()

// List returns a sorted slice of all available component names.
func List() []string {
	names := make([]string, 0, len(Registry))
	for name := range Registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Get returns the raw .kiw template content for the named component.
// Name can be provided as "Button" or "Button.kiw".
func Get(name string) ([]byte, error) {
	name = strings.TrimSuffix(name, ".kiw")
	filename := name + ".kiw"
	body, err := FS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("component %q not found in Krewire UI: %w", name, err)
	}
	return body, nil
}
