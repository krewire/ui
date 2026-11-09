// Command manifestgen refreshes the derived fields of components.json
// (file, props, defaults, dependencies, styles, script) from the .kiw sources
// and fills a default preview for entries that have none.
//
// Hand-authored fields (name, category, description, preview, packs) are
// kept. Run it from the components directory:
//
//	go generate ./components
//
// With -check it exits non-zero instead of writing when the manifest is stale.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/krewire/ui/components"
	"github.com/krewire/ui/components/internal/analyze"
)

func main() {
	dir := flag.String("dir", ".", "components directory")
	check := flag.Bool("check", false, "report a stale manifest instead of rewriting it")
	flag.Parse()
	if err := run(*dir, *check); err != nil {
		fmt.Fprintln(os.Stderr, "manifestgen:", err)
		os.Exit(1)
	}
}

func run(dir string, check bool) error {
	path := filepath.Join(dir, components.ManifestFile)
	old, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var m components.Manifest
	if err := json.Unmarshal(old, &m); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	authored := map[string]components.ComponentMeta{}
	for _, c := range m.Components {
		authored[c.Name] = c
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.kiw"))
	if err != nil {
		return err
	}
	m.Components = m.Components[:0]
	for _, f := range files {
		src, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.Base(f), ".kiw")
		c, ok := authored[name]
		if !ok {
			return fmt.Errorf("%s has no entry in %s: add name, category and description", filepath.Base(f), components.ManifestFile)
		}
		m.Components = append(m.Components, analyze.Derive(c, string(src)))
	}
	sort.Slice(m.Components, func(i, j int) bool { return m.Components[i].Name < m.Components[j].Name })

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(m); err != nil {
		return err
	}
	if bytes.Equal(buf.Bytes(), old) {
		return nil
	}
	if check {
		return fmt.Errorf("%s is stale: run go generate ./components", path)
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
