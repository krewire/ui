package kiw

import (
	"strings"
	"unicode"
)

// initialisms are prop names that Go-style templates spell in capitals
// (.ID, .URL); a lower-case attribute (id="x") also resolves to them.
var initialisms = map[string]bool{"id": true, "url": true, "uri": true, "api": true, "html": true, "css": true, "svg": true}

// PropAliases returns the keys under which a component prop passed as name is
// visible to the component template, implementing the lenient prop-casing
// convention:
//
//   - the name as written: variant, Variant, property-name
//   - first letter upper and lower: Variant <-> variant, so {variant} and
//     .Variant read the same prop
//   - kebab-case and snake_case also as camelCase and PascalCase:
//     aria-label -> ariaLabel, AriaLabel
//     property-name -> propertyName, PropertyName
//   - common initialisms in capitals: id -> ID, url -> URL, github-url -> gitHubURL
//
// The name as written is always first. Earlier keys win when two props alias
// each other, so an explicit Variant= is never overwritten by variant=.
func PropAliases(name string) []string {
	if name == "" {
		return nil
	}
	out := []string{name}
	add := func(s string) {
		for _, v := range out {
			if v == s {
				return
			}
		}
		out = append(out, s)
	}
	base := name
	if strings.ContainsAny(name, "-_") {
		parts := strings.FieldsFunc(name, func(r rune) bool { return r == '-' || r == '_' })
		for i := 1; i < len(parts); i++ {
			parts[i] = upperFirst(parts[i])
		}
		base = strings.Join(parts, "")
	}
	add(upperFirst(base))
	add(lowerFirstRune(base))
	if initialisms[strings.ToLower(base)] {
		add(strings.ToUpper(base))
	}

	// Initialisms and brand names for multi-part names (e.g. github-url -> gitHubURL, githubURL)
	if strings.ContainsAny(name, "-_") {
		parts := strings.FieldsFunc(name, func(r rune) bool { return r == '-' || r == '_' })
		hasInitialism := false
		for i := 0; i < len(parts); i++ {
			pLower := strings.ToLower(parts[i])
			if initialisms[pLower] || pLower == "github" {
				hasInitialism = true
				break
			}
		}
		if hasInitialism {
			var initCamel, initPascal []string
			for i, p := range parts {
				pLower := strings.ToLower(p)
				var capPart string
				if initialisms[pLower] {
					capPart = strings.ToUpper(p)
				} else if pLower == "github" {
					if i == 0 {
						capPart = "gitHub"
					} else {
						capPart = "GitHub"
					}
				} else {
					capPart = upperFirst(p)
				}
				if i == 0 {
					if pLower == "github" {
						initCamel = append(initCamel, "gitHub")
						initPascal = append(initPascal, "GitHub")
					} else {
						initCamel = append(initCamel, lowerFirstRune(capPart))
						initPascal = append(initPascal, upperFirst(capPart))
					}
				} else {
					initCamel = append(initCamel, capPart)
					initPascal = append(initPascal, capPart)
				}
			}
			add(strings.Join(initPascal, ""))
			add(strings.Join(initCamel, ""))
			if strings.Contains(strings.ToLower(name), "github") && strings.Contains(strings.ToLower(name), "url") {
				add("gitHubURL")
				add("githubURL")
				add("GitHubURL")
			}
			if strings.Contains(strings.ToLower(name), "hide") && strings.Contains(strings.ToLower(name), "github") {
				add("hideGitHub")
				add("HideGitHub")
			}
		}
	}
	return out
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func lowerFirstRune(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
