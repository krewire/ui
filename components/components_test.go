package components_test

import (
	"strings"
	"testing"

	"github.com/krewire/ui/components"
)

func TestEveryEmbeddedComponentIsReadable(t *testing.T) {
	entries, err := components.FS.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	var count int
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".kiw") {
			continue
		}
		count++
		body, err := components.FS.ReadFile(e.Name())
		if err != nil {
			t.Errorf("%s: read: %v", e.Name(), err)
			continue
		}
		if strings.TrimSpace(string(body)) == "" {
			t.Errorf("%s: component is empty", e.Name())
		}

		name := strings.TrimSuffix(e.Name(), ".kiw")
		if _, exists := components.Registry[name]; !exists {
			t.Errorf("%s: exists in embedded FS but is missing from Registry", e.Name())
		}
	}

	if count != 66 {
		t.Fatalf("expected 66 .kiw components, got %d", count)
	}
}

func TestListAndGet(t *testing.T) {
	list := components.List()
	if len(list) != 66 {
		t.Fatalf("expected 66 components in List(), got %d", len(list))
	}

	body, err := components.Get("Button")
	if err != nil {
		t.Fatalf("Get(Button) failed: %v", err)
	}
	if !strings.Contains(string(body), "button") {
		t.Errorf("expected Button.kiw to contain button tag/class, got: %s", string(body)[:100])
	}

	// Should also work with .kiw extension
	body2, err := components.Get("HStack.kiw")
	if err != nil {
		t.Fatalf("Get(HStack.kiw) failed: %v", err)
	}
	if len(body2) == 0 {
		t.Errorf("expected non-empty HStack.kiw")
	}

	// Test newly introduced components
	newComponents := []string{
		"Toast", "Skeleton", "Banner", "Empty", "Stat",
		"Timeline", "Accordion", "Tag", "Pagination", "Popover",
		"Dropdown", "SearchInput", "Slider", "ButtonGroup", "ActionMenu",
	}

	for _, name := range newComponents {
		data, err := components.Get(name)
		if err != nil {
			t.Errorf("Get(%s) failed: %v", name, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("expected non-empty body for %s.kiw", name)
		}
	}

	// Non-existent component
	_, err = components.Get("NonExistent")
	if err == nil {
		t.Errorf("expected error for non-existent component, got nil")
	}
}
