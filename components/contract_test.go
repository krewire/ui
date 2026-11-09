package components_test

// Contract tests for the kui <-> .kiw integration boundary documented in
// ARCHITECTURE.md. They pin down what a .kiw renderer must provide and what
// kui is allowed to depend on, so drift fails here instead of at render time.

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/krewire/ui/components"
)

// rendererHelpers are the non-builtin template functions kui components may
// call. A .kiw renderer must provide them (ARCHITECTURE.md, "Renderer
// requirements").
var rendererHelpers = map[string]bool{
	"class":      true,
	"assetLinks": true,
}

// goTemplateWords are text/template keywords and builtin functions that every
// renderer built on html/template provides.
var goTemplateWords = map[string]bool{
	"if": true, "else": true, "end": true, "range": true, "with": true,
	"eq": true, "ne": true, "lt": true, "le": true, "gt": true, "ge": true,
	"and": true, "or": true, "not": true, "len": true, "index": true,
	"print": true, "printf": true, "println": true,
	"nil": true, "true": true, "false": true,
}

var categories = map[string]bool{
	"Layout": true, "Atoms": true, "Forms": true, "Navigation": true,
	"Feedback": true, "Overlays": true, "Data Display": true, "Document": true,
}

var (
	styleBlockRe   = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	scriptBlockRe  = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	htmlCommentRe  = regexp.MustCompile(`(?s)<!--.*?-->`)
	goActionRe     = regexp.MustCompile(`(?s)\{\{-?(.*?)-?\}\}`)
	stringLitRe    = regexp.MustCompile(`"(?:[^"\\]|\\.)*"|` + "`[^`]*`")
	identRe        = regexp.MustCompile(`[.$]?[A-Za-z_][A-Za-z0-9_]*`)
	componentTagRe = regexp.MustCompile(`<([A-Z][A-Za-z0-9]*)[\s/>]`)
)

// markup returns the template part of a component: everything except its
// <style> and <script> blocks and HTML comments.
func markup(src string) string {
	src = styleBlockRe.ReplaceAllString(src, "")
	src = scriptBlockRe.ReplaceAllString(src, "")
	return htmlCommentRe.ReplaceAllString(src, "")
}

func eachComponent(t *testing.T, fn func(name, src string)) {
	t.Helper()
	matches, err := fs.Glob(components.FS, "*.kiw")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range matches {
		body, err := components.FS.ReadFile(file)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		fn(strings.TrimSuffix(file, ".kiw"), string(body))
	}
}

func TestRegistryAndFSAgree(t *testing.T) {
	for key, meta := range components.Registry {
		if meta.Name != key {
			t.Errorf("Registry[%q].Name = %q", key, meta.Name)
		}
		if !categories[meta.Category] {
			t.Errorf("Registry[%q].Category = %q is not a documented category", key, meta.Category)
		}
		if strings.TrimSpace(meta.Description) == "" {
			t.Errorf("Registry[%q] has no description", key)
		}
		if _, err := components.FS.ReadFile(key + ".kiw"); err != nil {
			t.Errorf("Registry[%q] has no embedded %s.kiw: %v", key, key, err)
		}
	}
}

func TestComponentsUseOnlyContractHelpers(t *testing.T) {
	eachComponent(t, func(name, src string) {
		for _, action := range goActionRe.FindAllStringSubmatch(markup(src), -1) {
			code := stringLitRe.ReplaceAllString(action[1], "")
			for _, word := range identRe.FindAllString(code, -1) {
				if strings.HasPrefix(word, ".") || strings.HasPrefix(word, "$") {
					continue // field or variable access
				}
				if !goTemplateWords[word] && !rendererHelpers[word] {
					t.Errorf("%s: template action uses %q, which is not a contract helper (%v)", name, word, action[0])
				}
			}
		}
	})
}

func TestNestedComponentsAreRegistered(t *testing.T) {
	eachComponent(t, func(name, src string) {
		for _, m := range componentTagRe.FindAllStringSubmatch(markup(src), -1) {
			if _, ok := components.Registry[m[1]]; !ok {
				t.Errorf("%s: nests <%s>, which is not a kui component", name, m[1])
			}
		}
	})
}

// TestModuleHasNoExternalDependencies enforces the integration boundary: kui
// must not import the kiw CLI, the krewire core module, or any third-party
// package. Only the standard library and kui's own packages are allowed.
func TestModuleHasNoExternalDependencies(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	gomod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(gomod), "module github.com/krewire/ui\n") {
		t.Fatalf("unexpected go.mod at %s", root)
	}
	for _, line := range strings.Split(string(gomod), "\n") {
		if f := strings.Fields(line); len(f) > 0 && (f[0] == "require" || f[0] == "replace") {
			t.Errorf("go.mod declares %q; kui must stay dependency-free", strings.TrimSpace(line))
		}
	}

	fset := token.NewFileSet()
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if name := d.Name(); path != root && (strings.HasPrefix(name, ".") || name == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		for _, imp := range f.Imports {
			p, _ := strconv.Unquote(imp.Path.Value)
			first := strings.SplitN(p, "/", 2)[0]
			stdlib := !strings.Contains(first, ".")
			own := p == "github.com/krewire/ui" || strings.HasPrefix(p, "github.com/krewire/ui/")
			if !stdlib && !own {
				t.Errorf("%s imports %q; kui may import only the standard library and its own packages", rel, p)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
