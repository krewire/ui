package styles

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	// Matches class="..." or class='...'
	classAttrRe = regexp.MustCompile(`class\s*=\s*["']([^"']+)["']`)
	// Matches {{class "btn" .Variant ...}}
	classHelperRe = regexp.MustCompile(`\{\{\s*class\s+([^}]+)\}\}`)
	// String literal matcher inside helpers
	strLitRe = regexp.MustCompile(`"([^"\\]*)"|'([^'\\]*)'`)
)

// ScanSource extracts all distinct CSS class candidates from source code.
func ScanSource(src string) []string {
	seen := make(map[string]bool)

	// 1. Scan standard HTML/JSX/Kiw class attributes
	for _, m := range classAttrRe.FindAllStringSubmatch(src, -1) {
		if len(m) > 1 {
			for _, cls := range strings.Fields(m[1]) {
				cls = cleanCandidate(cls)
				if cls != "" {
					seen[cls] = true
				}
			}
		}
	}

	// 2. Scan Kiw template {{class "base" ...}} helpers
	for _, m := range classHelperRe.FindAllStringSubmatch(src, -1) {
		if len(m) > 1 {
			for _, sm := range strLitRe.FindAllStringSubmatch(m[1], -1) {
				val := sm[1]
				if val == "" {
					val = sm[2]
				}
				for _, cls := range strings.Fields(val) {
					cls = cleanCandidate(cls)
					if cls != "" {
						seen[cls] = true
					}
				}
			}
		}
	}

	// 3. Scan data-kiw-component attributes
	compRe := regexp.MustCompile(`data-kiw-component\s*=\s*["']([^"']+)["']`)
	for _, m := range compRe.FindAllStringSubmatch(src, -1) {
		if len(m) > 1 {
			compName := strings.ToLower(m[1])
			seen[compName] = true
		}
	}

	res := make([]string, 0, len(seen))
	for cls := range seen {
		res = append(res, cls)
	}
	return res
}

func cleanCandidate(cls string) string {
	cls = strings.TrimSpace(cls)
	cls = strings.Trim(cls, `{}()[]"';:,`)
	if strings.HasPrefix(cls, ".") {
		cls = strings.TrimPrefix(cls, ".")
	}
	// Validate characters: must be valid CSS class character sequence
	for _, r := range cls {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' && r != ':' && r != '/' && r != '.' {
			return ""
		}
	}
	return cls
}

// ScanDir walks root directory and extracts all distinct CSS class candidates from
// .kiw, .html, .go, and .md files.
func ScanDir(root string) ([]string, error) {
	seen := make(map[string]bool)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".kiw" || ext == ".html" || ext == ".go" || ext == ".md" {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil // skip unreadable
			}
			for _, cls := range ScanSource(string(data)) {
				seen[cls] = true
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(seen))
	for cls := range seen {
		res = append(res, cls)
	}
	return res, nil
}
