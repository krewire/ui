package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ManifestFile is the name of the embedded component manifest. Its structure
// is described by components.schema.json (JSON Schema Draft-07).
const ManifestFile = "components.json"

// ComponentMeta is the catalog entry of one .kiw component.
//
// Name, Category, Description and Preview are authored by hand in
// components.json. Props, Defaults, Dependencies, Styles and Script are
// derived from the component source by `go generate ./components` and are
// verified by the contract tests, so they never drift from the markup.
type ComponentMeta struct {
	// Name is the component name and the file name without .kiw.
	Name string `json:"name"`
	// Category is one of the documented catalog categories.
	Category string `json:"category"`
	// Description is a one-line summary.
	Description string `json:"description"`
	// File is the component's file name inside FS.
	File string `json:"file"`
	// Props are the props and slots the markup reads, in their documented
	// lower-camel spelling. "body" is the child content of a paired tag.
	Props []string `json:"props,omitempty"`
	// Defaults are the values the class helper falls back to (variant, size, ...).
	Defaults map[string]string `json:"defaults,omitempty"`
	// Dependencies are the other kui components the markup nests.
	Dependencies []string `json:"dependencies,omitempty"`
	// Styles reports whether the component ships a scoped <style> block.
	Styles bool `json:"styles"`
	// Script reports whether the component ships an inline <script> block.
	Script bool `json:"script"`
	// Preview is a minimal usage example in .kiw markup.
	Preview string `json:"preview"`
}

// Pack is a curated starter bundle of components installed together.
type Pack struct {
	// Name is the pack identifier used on the command line (kiw add ui:<name>).
	Name string `json:"name"`
	// Description is a one-line summary.
	Description string `json:"description"`
	// Components lists the pack's components; dependencies are added on install.
	Components []string `json:"components"`
}

// Manifest is the decoded components.json catalog.
type Manifest struct {
	Schema     string          `json:"$schema,omitempty"`
	Version    int             `json:"version"`
	Packs      []Pack          `json:"packs"`
	Components []ComponentMeta `json:"components"`
}

// LoadManifest decodes the embedded components.json. Unknown fields are
// rejected so the manifest cannot silently drift from its schema.
func LoadManifest() (*Manifest, error) {
	data, err := FS.ReadFile(ManifestFile)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var m Manifest
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("components: decode %s: %w", ManifestFile, err)
	}
	return &m, nil
}

func mustRegistry() map[string]ComponentMeta {
	m, err := LoadManifest()
	if err != nil {
		panic(err)
	}
	reg := make(map[string]ComponentMeta, len(m.Components))
	for _, c := range m.Components {
		reg[c.Name] = c
	}
	return reg
}

// Packs returns the curated starter packs, sorted by name.
func Packs() []Pack {
	m, err := LoadManifest()
	if err != nil {
		return nil
	}
	sort.Slice(m.Packs, func(i, j int) bool { return m.Packs[i].Name < m.Packs[j].Name })
	return m.Packs
}

// Lookup resolves a component name case-insensitively ("button", "Button",
// "Button.kiw") to its catalog entry.
func Lookup(name string) (ComponentMeta, bool) {
	name = strings.TrimSuffix(name, ".kiw")
	if c, ok := Registry[name]; ok {
		return c, true
	}
	for key, c := range Registry {
		if strings.EqualFold(key, name) {
			return c, true
		}
	}
	return ComponentMeta{}, false
}

// Resolve expands component names and "pack:<name>" references into the
// sorted set of components to install, including transitive dependencies.
func Resolve(refs ...string) ([]string, error) {
	packs := map[string]Pack{}
	for _, p := range Packs() {
		packs[p.Name] = p
	}
	seen := map[string]bool{}
	var visit func(string) error
	visit = func(ref string) error {
		if name, ok := strings.CutPrefix(ref, "pack:"); ok {
			p, ok := packs[strings.ToLower(name)]
			if !ok {
				return fmt.Errorf("components: unknown pack %q", name)
			}
			for _, c := range p.Components {
				if err := visit(c); err != nil {
					return err
				}
			}
			return nil
		}
		c, ok := Lookup(ref)
		if !ok {
			return fmt.Errorf("components: unknown component %q", ref)
		}
		if seen[c.Name] {
			return nil
		}
		seen[c.Name] = true
		for _, d := range c.Dependencies {
			if err := visit(d); err != nil {
				return err
			}
		}
		return nil
	}
	for _, r := range refs {
		if err := visit(r); err != nil {
			return nil, err
		}
	}
	out := make([]string, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	sort.Strings(out)
	return out, nil
}
