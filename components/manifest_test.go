package components_test

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/krewire/ui/components"
)

type schemaDoc struct {
	Required    []string `json:"required"`
	Definitions struct {
		Component struct {
			Required   []string `json:"required"`
			Properties struct {
				Category struct {
					Enum []string `json:"enum"`
				} `json:"category"`
			} `json:"properties"`
		} `json:"component"`
		Pack struct {
			Properties struct {
				Name struct {
					Pattern string `json:"pattern"`
				} `json:"name"`
			} `json:"properties"`
		} `json:"pack"`
	} `json:"definitions"`
}

func loadSchema(t *testing.T) schemaDoc {
	t.Helper()
	data, err := components.FS.ReadFile(components.ManifestFile)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if raw["$schema"] != "./components.schema.json" {
		t.Errorf("components.json $schema = %v, want ./components.schema.json", raw["$schema"])
	}
	schemaBytes, err := os.ReadFile("components.schema.json") // tests run in components/
	if err != nil {
		t.Fatal(err)
	}
	var s schemaDoc
	if err := json.Unmarshal(schemaBytes, &s); err != nil {
		t.Fatalf("components.schema.json: %v", err)
	}
	return s
}

func TestManifestConformsToSchema(t *testing.T) {
	s := loadSchema(t)
	var enum []string
	for c := range categories {
		enum = append(enum, c)
	}
	sort.Strings(enum)
	schemaEnum := append([]string(nil), s.Definitions.Component.Properties.Category.Enum...)
	sort.Strings(schemaEnum)
	if !reflect.DeepEqual(enum, schemaEnum) {
		t.Errorf("schema category enum %v != documented categories %v", schemaEnum, enum)
	}

	m, err := components.LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	nameRe := regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	for _, c := range m.Components {
		fields := map[string]string{"name": c.Name, "category": c.Category, "description": c.Description, "file": c.File, "preview": c.Preview}
		for _, req := range s.Definitions.Component.Required {
			if v, ok := fields[req]; ok && strings.TrimSpace(v) == "" {
				t.Errorf("%s: required field %q is empty", c.Name, req)
			}
		}
		if !nameRe.MatchString(c.Name) {
			t.Errorf("%s: name is not PascalCase", c.Name)
		}
		for _, d := range c.Dependencies {
			if _, ok := components.Registry[d]; !ok {
				t.Errorf("%s: dependency %q is not a kui component", c.Name, d)
			}
		}
	}

	packRe := regexp.MustCompile(s.Definitions.Pack.Properties.Name.Pattern)
	seen := map[string]bool{}
	for _, p := range m.Packs {
		if !packRe.MatchString(p.Name) || seen[p.Name] {
			t.Errorf("pack %q: invalid or duplicate name", p.Name)
		}
		seen[p.Name] = true
		for _, c := range p.Components {
			if _, ok := components.Registry[c]; !ok {
				t.Errorf("pack %q: %q is not a kui component", p.Name, c)
			}
		}
	}
	for _, want := range []string{"base", "minimal", "dashboard"} {
		if !seen[want] {
			t.Errorf("starter pack %q is missing", want)
		}
	}
}

func TestResolve(t *testing.T) {
	got, err := components.Resolve("pack:minimal")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"Button", "Card", "Input"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Resolve(pack:minimal) = %v, want %v", got, want)
	}

	got, err = components.Resolve("navbar")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Navbar", "Brand", "Hamburger", "NavMenu", "ThemeSwitch"} {
		if !contains(got, want) {
			t.Errorf("Resolve(navbar) = %v, missing dependency %s", got, want)
		}
	}

	if _, err := components.Resolve("Nope"); err == nil {
		t.Error("Resolve(Nope): want error")
	}
	if _, err := components.Resolve("pack:nope"); err == nil {
		t.Error("Resolve(pack:nope): want error")
	}
	if c, ok := components.Lookup("button.kiw"); !ok || c.Name != "Button" {
		t.Errorf("Lookup(button.kiw) = %v, %v", c.Name, ok)
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
