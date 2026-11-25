package widget

import (
	"strings"
	"testing"
)

// TestSafeURL_BlocksExecutableSchemes verifies that a link target
// cannot carry a script-executing scheme. HTML-escaping alone does not help
// here: "javascript:alert(1)" contains no character that needs escaping, so it
// survives attribute encoding and still runs when the link is clicked.
func TestSafeURL_BlocksExecutableSchemes(t *testing.T) {
	unsafe := []string{
		"javascript:alert(1)",
		"JaVaScRiPt:alert(1)",
		"  javascript:alert(1)",
		"java\tscript:alert(1)",
		"vbscript:msgbox(1)",
		"data:text/html,<script>alert(1)</script>",
		"file:///etc/passwd",
	}
	for _, raw := range unsafe {
		if got := SafeURL(raw); got != "#" {
			t.Errorf("SafeURL(%q) = %q, want %q", raw, got, "#")
		}
	}
}

// TestSafeURL_AllowsLegitimateSchemes guards against the
// validator breaking ordinary navigation.
func TestSafeURL_AllowsLegitimateSchemes(t *testing.T) {
	cases := map[string]string{
		"/dashboard":     "/dashboard",
		"../relative":    "../relative",
		"#section":       "#section",
		"?q=1":           "?q=1",
		"https://ok.ex":  "https://ok.ex",
		"http://ok.ex":   "http://ok.ex",
		"mailto:a@b.ex":  "mailto:a@b.ex",
		"//cdn.ex/x.png": "//cdn.ex/x.png",
	}
	for raw, want := range cases {
		if got := SafeURL(raw); got != want {
			t.Errorf("SafeURL(%q) = %q, want %q", raw, got, want)
		}
	}
}

// TestSafeURL_EscapesQuoteInjection verifies that a URL
// containing a double quote cannot break out of the href attribute.
func TestSafeURL_EscapesQuoteInjection(t *testing.T) {
	got := SafeURL(`https://ok.example/" onmouseover="alert(1)`)
	if strings.Contains(got, `" onmouseover=`) {
		t.Errorf("SafeURL must escape quotes, got %q", got)
	}
	if !strings.Contains(got, "&#34;") {
		t.Errorf("SafeURL must entity-encode quotes, got %q", got)
	}
}

// TestButtonLinkUsesSafeURL verifies the button link path routes
// through SafeURL rather than raw escaping.
func TestButtonLinkUsesSafeURL(t *testing.T) {
	out := string(NewButton("click").Link("javascript:alert(1)").Render())
	if strings.Contains(out, "javascript:") {
		t.Errorf("button link must not carry a javascript: URL:\n%s", out)
	}
}

func TestSanitizeHTML_StripsDangerousScriptsAndSinks(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		disallow []string
		mustHave []string
	}{
		{
			name:     "Script tag with body",
			input:    `<p>Hello <script>alert("pwned")</script>World</p>`,
			disallow: []string{"<script", "alert", "pwned", "</script>"},
			mustHave: []string{"<p>", "Hello World", "</p>"},
		},
		{
			name:     "Style tag with body",
			input:    `<style>body { display: none; }</style><p>Content</p>`,
			disallow: []string{"<style", "display: none", "</style>"},
			mustHave: []string{"<p>Content</p>"},
		},
		{
			name:     "Inline onerror on image",
			input:    `<img src="invalid.png" onerror="alert(document.cookie)" alt="Test">`,
			disallow: []string{"onerror", "alert", "document.cookie"},
			mustHave: []string{"<img", `src="invalid.png"`, `alt="Test"`},
		},
		{
			name:     "SVG with inline onload",
			input:    `<svg onload="alert(1)"><circle cx="50" cy="50" r="40"/></svg>`,
			disallow: []string{"<svg", "onload", "alert(1)", "<circle"},
			mustHave: []string{},
		},
		{
			name:     "Iframe sink",
			input:    `<iframe src="https://evil.com"></iframe><b>Safe</b>`,
			disallow: []string{"<iframe", "evil.com", "</iframe>"},
			mustHave: []string{"<b>Safe</b>"},
		},
		{
			name:     "Anchor with javascript URI",
			input:    `<a href="javascript:alert(1)">Click here</a>`,
			disallow: []string{"javascript:", "alert(1)"},
			mustHave: []string{`<a href="#">Click here</a>`},
		},
		{
			name:     "Anchor with target blank gets noopener",
			input:    `<a href="https://krewire.dev" target="_blank">Portal</a>`,
			disallow: []string{},
			mustHave: []string{`href="https://krewire.dev"`, `target="_blank"`, `rel="noopener noreferrer"`},
		},
		{
			name:     "Malformed unclosed tag",
			input:    `<b onmouseover="alert(1)"</b>Text`,
			disallow: []string{"onmouseover", "alert(1)"},
			mustHave: []string{"Text"},
		},
		{
			name:     "HTML comment stripped",
			input:    `<!-- secret comment --><b>Visible</b>`,
			disallow: []string{"secret comment"},
			mustHave: []string{"<b>Visible</b>"},
		},
		{
			name:     "Null byte stripped",
			input:    "<p>Null\x00Byte</p>",
			disallow: []string{"\x00"},
			mustHave: []string{"<p>NullByte</p>"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			sanitized := string(SanitizeHTML(tc.input))
			for _, bad := range tc.disallow {
				if strings.Contains(sanitized, bad) {
					t.Errorf("SanitizeHTML(%q) contains disallowed %q in %q", tc.input, bad, sanitized)
				}
			}
			for _, good := range tc.mustHave {
				if !strings.Contains(sanitized, good) {
					t.Errorf("SanitizeHTML(%q) missing expected %q in %q", tc.input, good, sanitized)
				}
			}
		})
	}
}

func TestButton_WithActionAndTrustedScript(t *testing.T) {
	// 1. WithAction (delegated event handler)
	btnAction := NewButton("Delete").WithAction("delete-item", "item_42")
	outAction := string(btnAction.Render())
	if !strings.Contains(outAction, `data-action="delete-item"`) {
		t.Errorf("expected data-action in output, got: %s", outAction)
	}
	if !strings.Contains(outAction, `data-action-payload="item_42"`) {
		t.Errorf("expected data-action-payload in output, got: %s", outAction)
	}

	// 2. WithTrustedClickHandler
	btnScript := NewButton("Run").WithTrustedClickHandler(TrustedScript("window.doSomething()"))
	outScript := string(btnScript.Render())
	if !strings.Contains(outScript, `onclick="window.doSomething()"`) {
		t.Errorf("expected onclick in output, got: %s", outScript)
	}

	// 3. TrustedHTML
	th := TrustedHTML("<span>Trusted Markup</span>")
	if string(th.Render()) != "<span>Trusted Markup</span>" {
		t.Errorf("expected raw markup preserved for TrustedHTML")
	}
}
