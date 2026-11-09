// Package analyze derives component metadata (props, defaults, nested
// component dependencies, style/script presence) from .kiw source.
//
// The derived fields of components.json are computed here so the manifest
// generator and the manifest contract test share one implementation.
package analyze

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/krewire/ui/components"
)

// Facts are the metadata derived from one component's source.
type Facts struct {
	Props        []string
	Defaults     map[string]string
	Dependencies []string
	Styles       bool
	Script       bool
}

var (
	styleBlockRe  = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
	scriptBlockRe = regexp.MustCompile(`(?is)<script([^>]*)>(.*?)</script>`)
	htmlCommentRe = regexp.MustCompile(`(?s)<!--.*?-->`)
	tagRe         = regexp.MustCompile(`<([A-Z][A-Za-z0-9]*)[\s/>]`)
	actionRe      = regexp.MustCompile(`(?s)\{\{-?(.*?)-?\}\}`)
	fieldRe       = regexp.MustCompile(`(?:^|[^\w$.])\.([A-Za-z_]\w*)`)
	stringLitRe   = regexp.MustCompile(`"(?:[^"\\]|\\.)*"|'[^']*'|` + "`[^`]*`")
	blockRe       = regexp.MustCompile(`\{(?:#if|:else\s+if)\s+([^}]*)\}`)
	eachRe        = regexp.MustCompile(`\{#each\s+([\w$.]+)\s+as\s+([\w$]+)(?:\s*,\s*([\w$]+))?\s*\}`)
	interpRe      = regexp.MustCompile(`\{\s*([A-Za-z_][\w$]*(?:\.[\w$]+)*)\s*\}`)
	identRe       = regexp.MustCompile(`[A-Za-z_$][\w$]*(?:\.[\w$]+)*`)
	classCallRe   = regexp.MustCompile(`\{\{-?\s*class\s+"[^"]*"((?:\s+\.[A-Za-z_]\w*\s+"[^"]*")*)`)
	classPairRe   = regexp.MustCompile(`\.([A-Za-z_]\w*)\s+"([^"]*)"`)
)

var literals = map[string]bool{"true": true, "false": true, "nil": true, "null": true}

// Markup returns the template part of a component: the source without its
// <style> and <script> blocks and HTML comments.
func Markup(src string) string {
	src = styleBlockRe.ReplaceAllString(src, "")
	src = scriptBlockRe.ReplaceAllString(src, "")
	return htmlCommentRe.ReplaceAllString(src, "")
}

// PropName normalizes a template reference to the documented prop spelling:
// the first path segment with a lower-case first letter (.Variant -> variant).
func PropName(ref string) string {
	ref = strings.TrimLeft(ref, ".")
	if i := strings.IndexByte(ref, '.'); i >= 0 {
		ref = ref[:i]
	}
	if ref == "" {
		return ""
	}
	// Lower the leading run of capitals, Go initialism style:
	// Variant -> variant, ID -> id, URLPath -> urlPath, GitHubURL -> gitHubURL.
	r := []rune(ref)
	n := 0
	for n < len(r) && unicode.IsUpper(r[n]) {
		n++
	}
	if n > 1 && n < len(r) {
		n-- // keep the capital that starts the next word
	}
	for i := 0; i < n || i == 0; i++ {
		r[i] = unicode.ToLower(r[i])
	}
	return string(r)
}

// Analyze derives the Facts of the component named name from its source.
func Analyze(name, src string) Facts {
	f := Facts{Defaults: map[string]string{}}
	f.Styles = styleBlockRe.MatchString(src)
	for _, m := range scriptBlockRe.FindAllStringSubmatch(src, -1) {
		if strings.TrimSpace(m[2]) != "" {
			f.Script = true
		}
	}
	mk := Markup(src)

	loopVars := map[string]bool{}
	props := map[string]bool{}
	add := func(ref string) {
		root := ref
		if i := strings.IndexByte(root, '.'); i >= 0 {
			root = root[:i]
		}
		if root == "" || strings.HasPrefix(root, "$") || literals[root] || loopVars[root] {
			return
		}
		if p := PropName(root); p != "" {
			props[p] = true
		}
	}
	for _, m := range eachRe.FindAllStringSubmatch(mk, -1) {
		loopVars[m[2]] = true
		if m[3] != "" {
			loopVars[m[3]] = true
		}
	}
	for _, m := range eachRe.FindAllStringSubmatch(mk, -1) {
		add(m[1])
	}
	for _, m := range blockRe.FindAllStringSubmatch(mk, -1) {
		for _, id := range identRe.FindAllString(stringLitRe.ReplaceAllString(m[1], ""), -1) {
			add(id)
		}
	}
	for _, m := range interpRe.FindAllStringSubmatchIndex(mk, -1) {
		if m[0] > 0 && mk[m[0]-1] == '{' {
			continue // part of a {{ ... }} action
		}
		add(mk[m[2]:m[3]])
	}
	for _, m := range actionRe.FindAllStringSubmatch(mk, -1) {
		code := stringLitRe.ReplaceAllString(m[1], "")
		for _, fm := range fieldRe.FindAllStringSubmatch(code, -1) {
			add(fm[1])
		}
	}
	for p := range props {
		f.Props = append(f.Props, p)
	}
	sort.Strings(f.Props)

	for _, m := range classCallRe.FindAllStringSubmatch(mk, -1) {
		for _, pm := range classPairRe.FindAllStringSubmatch(m[1], -1) {
			f.Defaults[PropName(pm[1])] = pm[2]
		}
	}

	deps := map[string]bool{}
	for _, m := range tagRe.FindAllStringSubmatch(mk, -1) {
		if m[1] != name {
			deps[m[1]] = true
		}
	}
	for d := range deps {
		f.Dependencies = append(f.Dependencies, d)
	}
	sort.Strings(f.Dependencies)
	return f
}

// Derive returns c with its derived manifest fields (file, props, defaults,
// dependencies, styles, script) recomputed from src, and a default preview
// when none is authored.
func Derive(c components.ComponentMeta, src string) components.ComponentMeta {
	f := Analyze(c.Name, src)
	c.File = c.Name + ".kiw"
	c.Props = f.Props
	c.Defaults = f.Defaults
	if len(c.Defaults) == 0 {
		c.Defaults = nil
	}
	c.Dependencies = f.Dependencies
	c.Styles = f.Styles
	c.Script = f.Script
	if strings.TrimSpace(c.Preview) == "" {
		c.Preview = "<" + c.Name + " />"
		for _, p := range f.Props {
			if p == "body" {
				c.Preview = "<" + c.Name + ">…</" + c.Name + ">"
			}
		}
	}
	return c
}
