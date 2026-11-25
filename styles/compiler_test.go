package styles_test

import (
	"strings"
	"testing"

	"github.com/krewire/ui/styles"
)

func TestScanner_ExtractClasses(t *testing.T) {
	sample := `<div class="flex items-center justify-between p-4 bg-surface text-fg">
		{{class "btn" .Variant "primary" .Size "md"}}
		<span class="badge badge-green text-xs font-semibold">Active</span>
		<button class="hover:bg-primary sm:flex-col dark:bg-bg" data-kiw-component="Button">Click</button>
	</div>`

	classes := styles.ScanSource(sample)
	classMap := make(map[string]bool)
	for _, c := range classes {
		classMap[c] = true
	}

	expected := []string{
		"flex", "items-center", "justify-between", "p-4", "bg-surface", "text-fg",
		"btn", "badge", "badge-green", "text-xs", "font-semibold",
		"hover:bg-primary", "sm:flex-col", "dark:bg-bg", "button",
	}

	for _, exp := range expected {
		if !classMap[exp] {
			t.Errorf("expected scanned class %q, not found in %v", exp, classes)
		}
	}
}

func TestCompiler_TreeShaking_OnDemand(t *testing.T) {
	classes := []string{
		"flex",
		"items-center",
		"p-4",
		"sm:flex",
		"btn",
		"btn-primary",
		"hover:bg-primary",
		"dark:bg-surface",
	}

	c := styles.NewCompiler()
	out := c.GenerateCSS(classes)

	// 1. Used classes must be present
	mustHave := []string{
		".flex { display: flex; }",
		".items-center { align-items: center; }",
		".p-4 { padding: 1rem; }",
		".btn {",
		".btn-primary",
		"@media (min-width: 640px)",
		".sm\\:flex { display: flex; }",
		".hover\\:bg-primary:hover { background-color: var(--primary); }",
		"[data-theme=\"dark\"] .dark\\:bg-surface",
	}

	for _, token := range mustHave {
		if !strings.Contains(out, token) {
			t.Errorf("expected output to contain %q, but was missing", token)
		}
	}

	// 2. Unused classes must NOT be present (tree-shaken!)
	mustNotHave := []string{
		".grid-cols-12",
		".p-24",
		".text-5xl",
		".col-span-8",
		".w-fit",
	}

	for _, token := range mustNotHave {
		if strings.Contains(out, token) {
			t.Errorf("expected unused class %q to be purged from output", token)
		}
	}
}

func TestCompiler_SemanticComponents(t *testing.T) {
	c := styles.NewCompiler()
	// Only btn is used
	out := c.GenerateCSS([]string{"btn", "btn-primary"})
	if !strings.Contains(out, ".btn {") {
		t.Error("expected .btn in output")
	}
	// Card was not used, should not be included
	if strings.Contains(out, ".card {") {
		t.Error("expected .card to be omitted when not used")
	}
}
