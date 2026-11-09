package styles_test

import (
	"strings"
	"testing"

	"github.com/krewire/ui/styles"
)

func TestStylesheetContainsCoreTokens(t *testing.T) {
	css := styles.CSS
	if len(css) == 0 {
		t.Fatal("styles.CSS is empty")
	}

	essentialTokens := []string{
		"--primary",
		"--secondary",
		"--accent",
		"--bg",
		"--fg",
		"--border",
		"--pop-border",
		"--pop-shadow",
		"--ui-font-sans",
		"--ui-font-mono",
		"--ui-space-4",
		"--ui-radius",
		"[data-theme=\"dark\"]",
		".card",
		".badge",
		".border-pop",
		"--kui-duration-normal",
		"kui-fade-in",
		"kui-pop-in",
		".animate-in",
		".animate-out",
		"[data-state=\"open\"]",
		"--ui-container-px-sm",
		"--ui-container-px-md",
		"--ui-container-px-lg",
		"--ui-section-py-sm",
		"--ui-section-py-md",
		"--ui-section-py-lg",
		".container-fluid",
		".section-py",
		".sm\\:p-4",
		".md\\:py-6",
		".lg\\:gap-8",
		".xl\\:m-auto",
	}

	for _, token := range essentialTokens {
		if !strings.Contains(css, token) {
			t.Errorf("styles.CSS missing expected token or selector: %s", token)
		}
	}
}
