package styles

import (
	"fmt"
	"sort"
	"strings"
)

// Compiler orchestrates scanning and on-demand CSS generation for KrewireCSS.
type Compiler struct {
	tokens     ThemeTokens
	includeAll bool
}

// Option configures the Compiler.
type Option func(*Compiler)

// WithAllSemantic includes all semantic components regardless of scan results.
func WithAllSemantic(all bool) Option {
	return func(c *Compiler) {
		c.includeAll = all
	}
}

// NewCompiler creates a new KrewireCSS compiler instance.
func NewCompiler(opts ...Option) *Compiler {
	c := &Compiler{
		tokens:     LightTheme,
		includeAll: false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// GenerateCSS compiles on-demand CSS for the provided list of classes.
func (c *Compiler) GenerateCSS(classes []string) string {
	var sb strings.Builder

	// 1. Base tokens & CSS custom properties
	sb.WriteString(c.renderTokens())
	sb.WriteString("\n")

	// 2. Base Resets
	sb.WriteString(c.renderResets())
	sb.WriteString("\n")

	// 3. Semantic Components
	usedMap := make(map[string]bool, len(classes))
	for _, cls := range classes {
		usedMap[strings.ToLower(cls)] = true
	}
	if c.includeAll {
		for _, comp := range SemanticComponents {
			usedMap[comp.Trigger] = true
		}
	}
	semanticCSS := RenderSemanticCSS(usedMap)
	if semanticCSS != "" {
		sb.WriteString("/* --- Semantic Components --- */\n")
		sb.WriteString(semanticCSS)
	}

	// 4. Utility Classes (Purged / JIT On-Demand)
	baseRules := make([]*Rule, 0)
	mediaRules := make(map[string][]*Rule)
	seenSelectors := make(map[string]bool)

	for _, cls := range classes {
		rule, ok := ResolveUtility(cls)
		if !ok {
			continue
		}
		if seenSelectors[rule.Selector] {
			continue
		}
		seenSelectors[rule.Selector] = true

		if rule.Media == "" {
			baseRules = append(baseRules, rule)
		} else {
			mediaRules[rule.Media] = append(mediaRules[rule.Media], rule)
		}
	}

	// Sort base rules deterministically
	sort.Slice(baseRules, func(i, j int) bool {
		return baseRules[i].Selector < baseRules[j].Selector
	})

	if len(baseRules) > 0 {
		sb.WriteString("/* --- Utility Classes (On-Demand) --- */\n")
		for _, r := range baseRules {
			sb.WriteString(fmt.Sprintf("%s { %s }\n", r.Selector, r.Declaration))
		}
		sb.WriteString("\n")
	}

	// Render media queries in standard order
	mediaOrder := []string{
		"min-width: 640px",
		"min-width: 768px",
		"min-width: 1024px",
		"min-width: 1280px",
		"min-width: 1536px",
	}

	for _, m := range mediaOrder {
		rules, ok := mediaRules[m]
		if !ok || len(rules) == 0 {
			continue
		}
		sort.Slice(rules, func(i, j int) bool {
			return rules[i].Selector < rules[j].Selector
		})

		sb.WriteString(fmt.Sprintf("@media (%s) {\n", m))
		for _, r := range rules {
			sb.WriteString(fmt.Sprintf("  %s { %s }\n", r.Selector, r.Declaration))
		}
		sb.WriteString("}\n\n")
	}

	// 5. Animations & Keyframes
	sb.WriteString(c.renderAnimations())

	return sb.String()
}

// CompileDir scans all matching source files in root directory and compiles the purged CSS.
func (c *Compiler) CompileDir(root string) (string, error) {
	classes, err := ScanDir(root)
	if err != nil {
		return "", err
	}
	return c.GenerateCSS(classes), nil
}

// CompileSource scans a single source string and compiles the purged CSS.
func (c *Compiler) CompileSource(src string) string {
	classes := ScanSource(src)
	return c.GenerateCSS(classes)
}

func (c *Compiler) renderTokens() string {
	return `:root {
  --ui-font-sans: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  --ui-font-display: 'Outfit', 'Inter', system-ui, sans-serif;
  --ui-font-mono: 'Fira Code', 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;

  --ui-radius-xs: 4px;
  --ui-radius-sm: 6px;
  --ui-radius: 8px;
  --ui-radius-md: 10px;
  --ui-radius-lg: 12px;
  --ui-radius-xl: 16px;
  --radius: var(--ui-radius);
  --radius-sm: var(--ui-radius-sm);
  --radius-lg: var(--ui-radius-lg);

  --shadow-sm: 0 1px 2px 0 rgba(11, 31, 59, 0.05);
  --shadow: 0 1px 3px 0 rgba(11, 31, 59, 0.08), 0 1px 2px -1px rgba(11, 31, 59, 0.05);
  --shadow-md: 0 4px 6px -1px rgba(11, 31, 59, 0.08), 0 2px 4px -2px rgba(11, 31, 59, 0.04);
  --shadow-lg: 0 10px 15px -3px rgba(11, 31, 59, 0.08), 0 4px 6px -4px rgba(11, 31, 59, 0.03);
  --shadow-xl: 0 20px 25px -5px rgba(11, 31, 59, 0.08), 0 8px 10px -6px rgba(11, 31, 59, 0.03);

  --ui-shadow-sm: var(--shadow-sm);
  --ui-shadow: var(--shadow);
  --ui-shadow-md: var(--shadow-md);
  --ui-shadow-lg: var(--shadow-lg);
  --ui-shadow-xl: var(--shadow-xl);

  --pop-border: 1px solid var(--border, #e2ded4);
  --pop-border-thin: 1px solid var(--border-subtle, #ede9e0);
  --pop-shadow: var(--shadow-sm);
  --pop-shadow-lg: var(--shadow-md);
  --pop-shadow-sm: var(--shadow-sm);
  --pop-shadow-hover: var(--shadow-md);

  --bg: #FAF8F4;
  --bg-2: #F2EDE4;
  --bg-3: #E8E2D5;
  --fg: #0B1F3B;
  --fg-muted: #57677d;
  --fg-subtle: #8fa2bc;
  --border: #e2ded4;
  --border-subtle: #ede9e0;
  --border-focus: #39D353;

  --primary: #39D353;
  --primary-content: #0B1F3B;
  --primary-hover: #2ecc50;
  --primary-subtle: rgba(57, 211, 83, 0.14);

  --secondary: #00D1C1;
  --secondary-content: #0B1F3B;
  --secondary-hover: #00b8aa;
  --secondary-subtle: rgba(0, 209, 193, 0.14);

  --accent: #FF3B2E;
  --accent-content: #ffffff;
  --accent-hover: #e0281b;
  --accent-subtle: rgba(255, 59, 46, 0.14);

  --success: #39D353;
  --info: #0284c7;
  --warning: #f59e0b;
  --error: #FF3B2E;

  --surface: #ffffff;
  --surface-muted: #F2EDE4;
  --accent-cyan: #00D1C1;
}

[data-theme="dark"],
html.dark,
.dark {
  --bg: #0B1F3B;
  --bg-2: #071527;
  --bg-3: #0f284a;
  --fg: #FAF8F4;
  --fg-muted: #9cb0cb;
  --fg-subtle: #657e9f;
  --border: #1c3a63;
  --border-subtle: #132c4d;
  --border-focus: #39D353;

  --primary: #39D353;
  --primary-content: #0B1F3B;
  --primary-hover: #5eed75;
  --primary-subtle: rgba(57, 211, 83, 0.18);

  --secondary: #00D1C1;
  --secondary-content: #0B1F3B;
  --secondary-hover: #33ddd0;
  --secondary-subtle: rgba(0, 209, 193, 0.18);

  --accent: #FF5C52;
  --accent-content: #ffffff;
  --accent-hover: #ff7870;
  --accent-subtle: rgba(255, 92, 82, 0.18);

  --success: #39D353;
  --info: #38bdf8;
  --warning: #fbbf24;
  --error: #f87171;

  --surface: #0f284a;
  --surface-muted: #071527;

  --pop-border: 1px solid #1c3a63;
  --pop-border-thin: 1px solid #132c4d;
  --pop-shadow: var(--ui-shadow-sm);
  --pop-shadow-lg: var(--ui-shadow-md);
  --pop-shadow-sm: var(--ui-shadow-sm);
  --pop-shadow-hover: var(--ui-shadow-md);
}`
}

func (c *Compiler) renderResets() string {
	return `*, *::before, *::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  background-color: var(--bg);
  color: var(--fg);
  font-family: var(--ui-font-sans);
  line-height: 1.5;
  -webkit-font-smoothing: antialiased;
}

a {
  color: inherit;
  text-decoration: none;
}

button, input, select, textarea {
  font-family: inherit;
  font-size: inherit;
  line-height: inherit;
  color: inherit;
  margin: 0;
}`
}

func (c *Compiler) renderAnimations() string {
	return `:root {
  --kui-duration-fast: 150ms;
  --kui-duration-normal: 200ms;
  --kui-duration-slow: 300ms;
  --kui-ease-standard: cubic-bezier(0.16, 1, 0.3, 1);
  --kui-ease-out: cubic-bezier(0, 0, 0.2, 1);
  --kui-ease-in: cubic-bezier(0.4, 0, 1, 1);
  --kui-ease-pop: var(--kui-ease-standard);
}

@keyframes kui-fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes kui-fade-out {
  from { opacity: 1; }
  to { opacity: 0; }
}

@keyframes kui-zoom-in {
  from { opacity: 0; transform: scale(0.97) translateY(4px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

@keyframes kui-zoom-out {
  from { opacity: 1; transform: scale(1) translateY(0); }
  to { opacity: 0; transform: scale(0.97) translateY(4px); }
}

@keyframes kui-pop-in {
  from { opacity: 0; transform: scale(0.97) translateY(4px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

@keyframes kui-pop-out {
  from { opacity: 1; transform: scale(1) translateY(0); }
  to { opacity: 0; transform: scale(0.97) translateY(4px); }
}

.animate-in {
  animation-duration: var(--kui-duration-normal);
  animation-timing-function: var(--kui-ease-standard);
  animation-fill-mode: both;
}

.animate-out {
  animation-duration: var(--kui-duration-fast);
  animation-timing-function: var(--kui-ease-in);
  animation-fill-mode: both;
}

[data-state="open"] {
  display: block;
}`
}
