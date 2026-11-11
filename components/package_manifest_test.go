package components_test

import (
	"encoding/json"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/krewire/ui/components"
)

// packageManifest reflects the structure of kiw.pkg.json.
type packageManifest struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	License      string            `json:"license"`
	Main         string            `json:"main"`
	Components   map[string]string `json:"components"`
	Styles       string            `json:"styles"`
	Dependencies map[string]string `json:"dependencies"`
}

func TestPackageManifestCompleteness(t *testing.T) {
	pkgPath := filepath.Join("..", "kiw.pkg.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatalf("failed to read kiw.pkg.json: %v", err)
	}

	var pkg packageManifest
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatalf("invalid kiw.pkg.json format: %v", err)
	}

	if pkg.Name != "@krewire/ui" {
		t.Errorf("expected package name '@krewire/ui', got %q", pkg.Name)
	}
	if pkg.Version == "" {
		t.Error("package version must not be empty")
	}
	if pkg.License != "MIT" {
		t.Errorf("expected MIT license, got %q", pkg.License)
	}
	if pkg.Styles == "" {
		t.Error("package styles bundle must be specified")
	}

	// Verify that style file exists
	stylePath := filepath.Join("..", pkg.Styles)
	if _, err := os.Stat(stylePath); err != nil {
		t.Errorf("styles file %q not found: %v", stylePath, err)
	}

	// Verify that every component in components.json is exported in kiw.pkg.json
	manifest, err := components.LoadManifest()
	if err != nil {
		t.Fatalf("failed to load components.json: %v", err)
	}

	for _, c := range manifest.Components {
		exportedPath, ok := pkg.Components[c.Name]
		if !ok {
			t.Errorf("component %q is in components.json but missing from kiw.pkg.json", c.Name)
			continue
		}
		// Verify on-disk file exists
		realPath := filepath.Join("..", exportedPath)
		if _, err := os.Stat(realPath); err != nil {
			t.Errorf("exported component %q file %q does not exist: %v", c.Name, realPath, err)
		}
	}

	if len(pkg.Components) != len(manifest.Components) {
		t.Errorf("kiw.pkg.json has %d components, components.json has %d", len(pkg.Components), len(manifest.Components))
	}
}

func TestNoNpmOrKrewireDependency(t *testing.T) {
	// Zero Node.js / NPM bloat: No package.json or node_modules
	for _, forbidden := range []string{"../package.json", "../package-lock.json", "../node_modules"} {
		if _, err := os.Stat(forbidden); err == nil {
			t.Errorf("found forbidden npm artifact: %s", forbidden)
		}
	}

	// Downward layering guard: ui package imports must never import github.com/krewire/krewire
	fset := token.NewFileSet()
	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) == ".go" && !strings.HasSuffix(path, "_test.go") {
			node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, imp := range node.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				if strings.HasPrefix(importPath, "github.com/krewire/krewire") {
					t.Errorf("%s imports %s, violating layering discipline", path, importPath)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
