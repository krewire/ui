package widget

import (
	"bytes"
	"html/template"
	"strings"
	"unicode"
)

// AllowedHTMLTags defines the permitted HTML element names in SanitizeHTML.
var AllowedHTMLTags = map[string]bool{
	"p": true, "br": true, "hr": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"b": true, "i": true, "strong": true, "em": true, "u": true, "s": true,
	"del": true, "ins": true, "sub": true, "sup": true,
	"code": true, "pre": true, "kbd": true, "samp": true,
	"blockquote": true, "q": true, "cite": true,
	"ul": true, "ol": true, "li": true, "dl": true, "dt": true, "dd": true,
	"span": true, "div": true,
	"a": true, "img": true,
	"table": true, "thead": true, "tbody": true, "tfoot": true,
	"tr": true, "th": true, "td": true, "caption": true,
}

// GlobalAllowedAttributes defines attributes permitted on any allowed tag.
var GlobalAllowedAttributes = map[string]bool{
	"class": true,
	"title": true,
	"dir":   true,
	"lang":  true,
	"id":    true,
}

// DisallowedTagsWithContent names tags whose entire content block (body) must be stripped.
var DisallowedTagsWithContent = map[string]bool{
	"script":   true,
	"style":    true,
	"iframe":   true,
	"object":   true,
	"embed":    true,
	"applet":   true,
	"noframes": true,
	"noscript": true,
}

// SanitizeHTML filters an untrusted HTML string and returns safe template.HTML markup.
// It enforces an allow-list of safe semantic HTML tags, strips script tags, style blocks,
// SVG payloads, inline event handlers (on*), unsafe schemes (javascript:, data:),
// and entity-escapes unpermitted elements.
func SanitizeHTML(input string) template.HTML {
	if strings.TrimSpace(input) == "" {
		return ""
	}

	// Remove null bytes and dangerous control characters
	cleaned := strings.Map(func(r rune) rune {
		if r == 0 || (r < 32 && r != '\t' && r != '\n' && r != '\r') || r == 0x7f {
			return -1
		}
		return r
	}, input)

	var buf bytes.Buffer
	i := 0
	n := len(cleaned)

	for i < n {
		if cleaned[i] == '<' {
			// Check for HTML comment <!-- ... -->
			if strings.HasPrefix(cleaned[i:], "<!--") {
				end := strings.Index(cleaned[i:], "-->")
				if end != -1 {
					i += end + 3
					continue
				}
				// Unclosed comment -> consume rest
				break
			}

			// Find closing '>'
			end := strings.IndexByte(cleaned[i:], '>')
			if end == -1 {
				// Malformed unclosed tag: escape '<' and continue
				buf.WriteString("&lt;")
				i++
				continue
			}

			tagContent := cleaned[i+1 : i+end]
			i += end + 1

			isClosing := false
			trimmed := strings.TrimSpace(tagContent)
			if strings.HasPrefix(trimmed, "/") {
				isClosing = true
				trimmed = strings.TrimSpace(trimmed[1:])
			}

			// Extract tag name
			var tagName string
			var rest string
			splitIdx := strings.IndexFunc(trimmed, unicode.IsSpace)
			if splitIdx != -1 {
				tagName = strings.ToLower(trimmed[:splitIdx])
				rest = trimmed[splitIdx:]
			} else {
				// Self-closing slash at end
				trimmed = strings.TrimSuffix(trimmed, "/")
				tagName = strings.ToLower(trimmed)
			}

			// If tag has content that must be stripped completely (e.g. <script>...</script>)
			if DisallowedTagsWithContent[tagName] && !isClosing {
				closingTag := "</" + tagName + ">"
				closeIdx := strings.Index(strings.ToLower(cleaned[i:]), closingTag)
				if closeIdx != -1 {
					i += closeIdx + len(closingTag)
				} else {
					// No closing tag found, consume rest
					break
				}
				continue
			}

			// If tag is not in allowed list, skip it completely
			if !AllowedHTMLTags[tagName] {
				continue
			}

			if isClosing {
				buf.WriteString("</")
				buf.WriteString(tagName)
				buf.WriteString(">")
				continue
			}

			// Parse and sanitize attributes
			attrs := sanitizeAttributes(tagName, rest)
			buf.WriteString("<")
			buf.WriteString(tagName)
			if attrs != "" {
				buf.WriteString(" ")
				buf.WriteString(attrs)
			}
			buf.WriteString(">")
		} else {
			buf.WriteByte(cleaned[i])
			i++
		}
	}

	return template.HTML(buf.String())
}

// sanitizeAttributes parses an attribute string and outputs sanitized key="value" pairs.
func sanitizeAttributes(tagName, raw string) string {
	raw = strings.TrimSuffix(strings.TrimSpace(raw), "/")
	if raw == "" {
		return ""
	}

	var pairs []string
	i := 0
	n := len(raw)

	for i < n {
		for i < n && unicode.IsSpace(rune(raw[i])) {
			i++
		}
		if i >= n {
			break
		}

		// Read attribute name
		nameStart := i
		for i < n && raw[i] != '=' && !unicode.IsSpace(rune(raw[i])) && raw[i] != '>' && raw[i] != '/' {
			i++
		}
		if nameStart == i {
			i++
			continue
		}
		attrName := strings.ToLower(raw[nameStart:i])

		// Skip whitespace after name
		for i < n && unicode.IsSpace(rune(raw[i])) {
			i++
		}

		attrVal := ""
		if i < n && raw[i] == '=' {
			i++ // skip '='
			for i < n && unicode.IsSpace(rune(raw[i])) {
				i++
			}
			if i < n {
				if raw[i] == '"' || raw[i] == '\'' {
					quote := raw[i]
					i++
					valStart := i
					for i < n && raw[i] != quote {
						i++
					}
					attrVal = raw[valStart:i]
					if i < n && raw[i] == quote {
						i++
					}
				} else {
					valStart := i
					for i < n && !unicode.IsSpace(rune(raw[i])) && raw[i] != '>' {
						i++
					}
					attrVal = raw[valStart:i]
				}
			}
		}

		// Security filtering:
		// 1. Disallow any on* event handlers (onclick, onerror, onload, etc.)
		if strings.HasPrefix(attrName, "on") {
			continue
		}
		// 2. Disallow style attribute (prevents CSS injection & expressions)
		if attrName == "style" {
			continue
		}

		// 3. Tag-specific attribute handling
		if tagName == "a" && attrName == "href" {
			safe := SafeURL(attrVal)
			pairs = append(pairs, `href="`+safe+`"`)
			continue
		}
		if tagName == "a" && attrName == "target" {
			low := strings.ToLower(attrVal)
			if low == "_blank" {
				pairs = append(pairs, `target="_blank" rel="noopener noreferrer"`)
				continue
			} else if low == "_self" {
				pairs = append(pairs, `target="_self"`)
				continue
			}
			continue
		}
		if tagName == "img" && attrName == "src" {
			safe := SafeURL(attrVal)
			pairs = append(pairs, `src="`+safe+`"`)
			continue
		}
		if (tagName == "th" || tagName == "td") && (attrName == "colspan" || attrName == "rowspan" || attrName == "scope") {
			pairs = append(pairs, attrName+`="`+template.HTMLEscapeString(attrVal)+`"`)
			continue
		}
		if tagName == "img" && (attrName == "alt" || attrName == "width" || attrName == "height") {
			pairs = append(pairs, attrName+`="`+template.HTMLEscapeString(attrVal)+`"`)
			continue
		}
		if GlobalAllowedAttributes[attrName] {
			pairs = append(pairs, attrName+`="`+template.HTMLEscapeString(attrVal)+`"`)
			continue
		}
	}

	return strings.Join(pairs, " ")
}
