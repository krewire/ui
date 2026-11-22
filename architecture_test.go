package ui_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// TestArchitecture_SubsystemLayering verifies downward-only and isolated
// dependencies across Krewire UI subsystem packages.
func TestArchitecture_SubsystemLayering(t *testing.T) {
	root := repoRoot(t)
	fset := token.NewFileSet()

	// Subsystem import invariants:
	// - Leaf packages must not import github.com/krewire/ui root facade.
	// - Specific leaf packages must remain strictly bounded.
	subsystemRules := map[string]struct {
		banned     []string
		mustBeLeaf bool
	}{
		"node": {
			banned:     []string{"github.com/krewire/ui"},
			mustBeLeaf: true,
		},
		"widget": {
			banned:     []string{"github.com/krewire/ui"},
			mustBeLeaf: true,
		},
		"form": {
			banned:     []string{"github.com/krewire/ui"},
			mustBeLeaf: true,
		},
		"styles": {
			banned:     []string{"github.com/krewire/ui"},
			mustBeLeaf: true,
		},
		"kiw": {
			banned:     []string{"github.com/krewire/ui"},
			mustBeLeaf: true,
		},
		"panel": {
			banned: []string{"github.com/krewire/ui", "github.com/krewire/ui/form", "github.com/krewire/ui/kiw"},
		},
		"dsl": {
			banned: []string{"github.com/krewire/ui", "github.com/krewire/ui/form", "github.com/krewire/ui/panel"},
		},
	}

	for pkgName, rule := range subsystemRules {
		pkgDir := filepath.Join(root, pkgName)
		if _, err := os.Stat(pkgDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return err
			}

			node, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
			if parseErr != nil {
				t.Fatalf("failed to parse %s: %v", path, parseErr)
			}

			for _, imp := range node.Imports {
				val, _ := strconv.Unquote(imp.Path.Value)

				// Check banned list
				for _, b := range rule.banned {
					if val == b {
						t.Errorf("layering violation in %s: %s must not import %s", path, pkgName, val)
					}
				}

				// If strictly leaf, must not import ANY package from github.com/krewire/ui
				if rule.mustBeLeaf {
					if strings.HasPrefix(val, "github.com/krewire/ui") {
						t.Errorf("leaf violation in %s: leaf package %s must not import internal package %s", path, pkgName, val)
					}
				}
			}
			return nil
		})

		if err != nil {
			t.Fatalf("failed walking %s: %v", pkgName, err)
		}
	}
}

// TestArchitecture_EcosystemIsolation verifies that the UI framework never
// imports higher layers (CLI toolchain, Cloud, Auth, Dash).
func TestArchitecture_EcosystemIsolation(t *testing.T) {
	root := repoRoot(t)
	fset := token.NewFileSet()

	forbiddenPrefixes := []string{
		"github.com/krewire/krew",
		"github.com/krewire/cloud",
		"github.com/krewire/auth",
		"github.com/krewire/dash",
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}

		// Skip vendor or git
		if strings.Contains(path, "/.git/") || strings.Contains(path, "/vendor/") {
			return nil
		}

		node, parseErr := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if parseErr != nil {
			t.Fatalf("failed to parse %s: %v", path, parseErr)
		}

		for _, imp := range node.Imports {
			val, _ := strconv.Unquote(imp.Path.Value)
			for _, forbidden := range forbiddenPrefixes {
				if strings.HasPrefix(val, forbidden) {
					t.Errorf("ecosystem isolation violation in %s: imports forbidden higher-layer package %s", path, val)
				}
			}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("failed checking ecosystem isolation: %v", err)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("cannot get working dir: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repo root with go.mod")
		}
		dir = parent
	}
}
