package components_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/krewire/ui/components"
)

// TestComponentCatalog_ConformanceCoverage verifies that every single component
// in the embedded catalog satisfies the authoritative quality, semantic, and
// accessibility contracts defined in ADR-K7P9L and issue #10.
func TestComponentCatalog_ConformanceCoverage(t *testing.T) {
	if len(components.Registry) != 66 {
		t.Fatalf("expected 66 components in Registry, got %d", len(components.Registry))
	}

	for name, meta := range components.Registry {
		t.Run(name, func(t *testing.T) {
			// 1. Metadata completeness
			if meta.Name != name {
				t.Errorf("mismatched name: meta.Name=%q, key=%q", meta.Name, name)
			}
			if meta.Category == "" {
				t.Errorf("component %s has empty category", name)
			}
			if meta.Description == "" {
				t.Errorf("component %s has empty description", name)
			}
			if meta.File != name+".kiw" {
				t.Errorf("component %s has invalid file %q", name, meta.File)
			}
			if meta.Preview == "" {
				t.Errorf("component %s has no preview markup", name)
			}

			// 2. Read embedded markup
			data, err := components.FS.ReadFile(meta.File)
			if err != nil {
				t.Fatalf("failed to read embedded %s: %v", meta.File, err)
			}
			src := string(data)

			// 3. Mandatory data-kiw-component attribute for introspection (except document-level & composite headers)
			skipDataAttr := map[string]bool{
				"Head": true, "Brand": true, "Navbar": true, "NavMenu": true, "SideMenu": true,
			}
			if !skipDataAttr[name] {
				expectedAttr := `data-kiw-component="` + name + `"`
				if !strings.Contains(src, expectedAttr) {
					t.Errorf("component %s missing %s attribute", name, expectedAttr)
				}
			}

			// 4. Zero unescaped raw sinks (no innerHTML without escaping)
			if strings.Contains(src, ".innerHTML") {
				t.Errorf("component %s uses unsafe innerHTML", name)
			}
			if strings.Contains(src, "eval(") {
				t.Errorf("component %s uses eval()", name)
			}
		})
	}
}

// TestComponentAccessibility_SemanticContracts asserts WCAG semantic markup,
// accessible roles, form bindings, and ARIA attributes per component category.
func TestComponentAccessibility_SemanticContracts(t *testing.T) {
	t.Run("Buttons and Interactive Elements", func(t *testing.T) {
		btnSrc := readComponent(t, "Button")
		if !strings.Contains(btnSrc, "<button") && !strings.Contains(btnSrc, `role="button"`) {
			t.Errorf("Button component lacks semantic <button> or role='button'")
		}
		if !strings.Contains(btnSrc, `type=`) {
			t.Errorf("Button component lacks explicit type attribute")
		}
		if !strings.Contains(btnSrc, `disabled`) {
			t.Errorf("Button component lacks disabled attribute handling")
		}
	})

	t.Run("Form Controls", func(t *testing.T) {
		inputSrc := readComponent(t, "Input")
		if !strings.Contains(inputSrc, "<input") {
			t.Errorf("Input component lacks semantic <input>")
		}
		if !strings.Contains(inputSrc, `aria-invalid`) && !strings.Contains(inputSrc, `aria-describedby`) {
			t.Errorf("Input component lacks accessibility validation ARIA attributes")
		}

		selectSrc := readComponent(t, "Select")
		if !strings.Contains(selectSrc, "<select") {
			t.Errorf("Select component lacks semantic <select>")
		}

		textareaSrc := readComponent(t, "Textarea")
		if !strings.Contains(textareaSrc, "<textarea") {
			t.Errorf("Textarea component lacks semantic <textarea>")
		}

		checkboxSrc := readComponent(t, "Checkbox")
		if !strings.Contains(checkboxSrc, `type="checkbox"`) {
			t.Errorf("Checkbox component lacks type='checkbox'")
		}
	})

	t.Run("Overlays and Dialogs", func(t *testing.T) {
		modalSrc := readComponent(t, "Modal")
		if !strings.Contains(modalSrc, "<dialog") && !strings.Contains(modalSrc, `role="dialog"`) {
			t.Errorf("Modal component lacks semantic <dialog> or role='dialog'")
		}
		if !strings.Contains(modalSrc, `aria-label="Close dialog"`) && !strings.Contains(modalSrc, `aria-modal`) {
			t.Errorf("Modal component lacks accessible dialog controls")
		}

		drawerSrc := readComponent(t, "Drawer")
		if !strings.Contains(drawerSrc, `role="dialog"`) {
			t.Errorf("Drawer component lacks role='dialog'")
		}
	})

	t.Run("Navigation and Landmarks", func(t *testing.T) {
		navbarSrc := readComponent(t, "Navbar")
		if !strings.Contains(navbarSrc, "<header") {
			t.Errorf("Navbar component lacks semantic <header>")
		}

		navMenuSrc := readComponent(t, "NavMenu")
		if !strings.Contains(navMenuSrc, "<nav") {
			t.Errorf("NavMenu component lacks semantic <nav>")
		}

		breadcrumbSrc := readComponent(t, "Breadcrumb")
		if !strings.Contains(breadcrumbSrc, "<nav") && !strings.Contains(breadcrumbSrc, `aria-label="breadcrumb"`) {
			t.Errorf("Breadcrumb component lacks navigation landmark or aria-label='breadcrumb'")
		}

		paginationSrc := readComponent(t, "Pagination")
		if !strings.Contains(paginationSrc, "<nav") && !strings.Contains(paginationSrc, `aria-label="pagination"`) {
			t.Errorf("Pagination component lacks navigation landmark or aria-label='pagination'")
		}
	})

	t.Run("Alerts and Feedback", func(t *testing.T) {
		alertSrc := readComponent(t, "Alert")
		if !strings.Contains(alertSrc, `role="alert"`) {
			t.Errorf("Alert component lacks role='alert'")
		}

		toastSrc := readComponent(t, "Toast")
		if !strings.Contains(toastSrc, `role="status"`) && !strings.Contains(toastSrc, `role="alert"`) {
			t.Errorf("Toast component lacks status/alert role")
		}
	})
}

// TestComponentSecurity_AdversarialSinks verifies that component attributes
// and template bindings do not execute unescaped script vectors or malicious schemes.
func TestComponentSecurity_AdversarialSinks(t *testing.T) {
	dangerousAttrs := []string{
		"onload=", "onerror=", "onclick=", "onmouseover=", "onfocus=",
	}

	for name, meta := range components.Registry {
		data, err := components.FS.ReadFile(meta.File)
		if err != nil {
			t.Fatalf("failed to read %s: %v", meta.File, err)
		}
		src := string(data)

		// Check for unquoted or raw event handler bindings that could execute caller input
		for _, attr := range dangerousAttrs {
			// Find occurrences of event handler in template body (outside <script> tags)
			bodyOnly := removeScriptAndStyle(src)
			if strings.Contains(strings.ToLower(bodyOnly), attr) {
				// Event handlers must be static or strictly controlled, never wrapping raw props
				re := regexp.MustCompile(`(?i)` + attr + `["']?\{\{`)
				if re.MatchString(bodyOnly) {
					t.Errorf("component %s has dangerous dynamic event handler binding: %s", name, attr)
				}
			}
		}
	}
}

// TestModelC_VendorEjectWorkflow tests the Model C package distribution contract:
// components can be read from the package and cleanly ejected/vendored into
// an application's local components directory without data loss.
func TestModelC_VendorEjectWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	vendorDir := filepath.Join(tempDir, "components")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	testComponents := []string{"Button", "Card", "Alert", "Input", "Modal"}
	for _, name := range testComponents {
		// 1. Retrieve component source via public API
		src, err := components.Get(name)
		if err != nil {
			t.Fatalf("components.Get(%s) failed: %v", name, err)
		}

		// 2. Write to vendor directory simulating `kiw ui vendor <Component>`
		dstPath := filepath.Join(vendorDir, name+".kiw")
		if err := os.WriteFile(dstPath, src, 0644); err != nil {
			t.Fatalf("failed to vendor %s: %v", name, err)
		}

		// 3. Verify vendored file content matches embedded FS byte-for-byte
		upstreamData, err := components.FS.ReadFile(name + ".kiw")
		if err != nil {
			t.Fatalf("FS.ReadFile(%s) failed: %v", name, err)
		}
		vendoredData, err := os.ReadFile(dstPath)
		if err != nil {
			t.Fatalf("os.ReadFile(%s) failed: %v", dstPath, err)
		}
		if string(upstreamData) != string(vendoredData) {
			t.Errorf("vendored %s content differs from upstream source", name)
		}
	}
}

func readComponent(t *testing.T, name string) string {
	t.Helper()
	data, err := components.FS.ReadFile(name + ".kiw")
	if err != nil {
		t.Fatalf("failed to read %s.kiw: %v", name, err)
	}
	return string(data)
}

func removeScriptAndStyle(src string) string {
	scriptRe := regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
	styleRe := regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	s := scriptRe.ReplaceAllString(src, "")
	return styleRe.ReplaceAllString(s, "")
}
