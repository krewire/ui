package kiw

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// YAMLUnmarshaler is an optional hook allowing external YAML decoders to be plugged in.
// If nil, the built-in pure-Go frontmatter parser is used.
var YAMLUnmarshaler func(data []byte, v any) error

// unmarshalYAML decodes YAML frontmatter into a destination pointer (typically *map[string]any).
func unmarshalYAML(data []byte, v any) error {
	if YAMLUnmarshaler != nil {
		return YAMLUnmarshaler(data, v)
	}
	out, ok := v.(*map[string]any)
	if !ok {
		return fmt.Errorf("yaml: unsupported target type %T", v)
	}
	parsed, err := parseYAMLFrontmatter(data)
	if err != nil {
		return err
	}
	*out = parsed
	return nil
}

type yamlLine struct {
	lineNum int
	indent  int
	raw     string
	trimmed string
}

// parseYAMLFrontmatter parses standard YAML frontmatter into a map[string]any.
// It is a zero-dependency, pure-Go parser designed specifically for .kiw frontmatter
// (key-value scalars, nested maps, and string lists).
func parseYAMLFrontmatter(data []byte) (map[string]any, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var lines []yamlLine
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Strip comments if not inside quotes
		trimmed = stripComment(trimmed)
		if trimmed == "" {
			continue
		}
		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		lines = append(lines, yamlLine{
			lineNum: lineNum,
			indent:  indent,
			raw:     raw,
			trimmed: trimmed,
		})
	}

	if len(lines) == 0 {
		return map[string]any{}, nil
	}

	idx := 0
	val, next, err := parseYAMLBlock(lines, 0, lines[0].indent)
	if err != nil {
		return nil, err
	}
	if next < len(lines) {
		return nil, fmt.Errorf("yaml: line %d: unexpected indentation or content", lines[next].lineNum)
	}
	res, ok := val.(map[string]any)
	if !ok {
		return map[string]any{}, nil
	}
	_ = idx
	return res, nil
}

func parseYAMLBlock(lines []yamlLine, start int, expectedIndent int) (any, int, error) {
	if start >= len(lines) {
		return nil, start, nil
	}

	if strings.HasPrefix(lines[start].trimmed, "- ") {
		// List block
		var list []any
		i := start
		for i < len(lines) {
			if lines[i].indent < expectedIndent {
				break
			}
			if lines[i].indent > expectedIndent {
				return nil, i, fmt.Errorf("yaml: line %d: unexpected list item indentation", lines[i].lineNum)
			}
			if !strings.HasPrefix(lines[i].trimmed, "- ") {
				break
			}
			itemText := strings.TrimSpace(lines[i].trimmed[2:])
			list = append(list, parseYAMLValue(itemText))
			i++
		}
		return list, i, nil
	}

	// Map block
	m := make(map[string]any)
	i := start
	var lastKey string
	var lastKeyHadValue bool

	for i < len(lines) {
		line := lines[i]
		if line.indent < expectedIndent {
			break
		}
		if line.indent > expectedIndent {
			if lastKeyHadValue {
				return nil, i, fmt.Errorf("yaml: line %d: mapping values are not allowed in this context", line.lineNum)
			}
			// Parse nested block
			nestedVal, next, err := parseYAMLBlock(lines, i, line.indent)
			if err != nil {
				return nil, next, err
			}
			m[lastKey] = nestedVal
			i = next
			lastKey = ""
			lastKeyHadValue = false
			continue
		}

		// Expected indent
		colonIdx := strings.IndexByte(line.trimmed, ':')
		if colonIdx < 0 {
			return nil, i, fmt.Errorf("yaml: line %d: expected key-value mapping", line.lineNum)
		}

		key := strings.TrimSpace(line.trimmed[:colonIdx])
		key = unquote(key)
		valStr := strings.TrimSpace(line.trimmed[colonIdx+1:])

		if valStr == "" {
			lastKey = key
			lastKeyHadValue = false
			m[key] = nil
			i++
		} else {
			lastKey = key
			lastKeyHadValue = true
			m[key] = parseYAMLValue(valStr)
			i++
		}
	}

	return m, i, nil
}

func stripComment(s string) string {
	inQuotes := false
	quoteChar := byte(0)
	for i := 0; i < len(s); i++ {
		b := s[i]
		if inQuotes {
			if b == quoteChar && (i == 0 || s[i-1] != '\\') {
				inQuotes = false
			}
			continue
		}
		if b == '"' || b == '\'' {
			inQuotes = true
			quoteChar = b
			continue
		}
		if b == '#' && (i == 0 || s[i-1] == ' ' || s[i-1] == '\t') {
			return strings.TrimSpace(s[:i])
		}
	}
	return s
}

func parseYAMLValue(s string) any {
	if s == "" || s == "~" || s == "null" {
		return nil
	}
	if s == "true" || s == "True" || s == "TRUE" {
		return true
	}
	if s == "false" || s == "False" || s == "FALSE" {
		return false
	}
	if (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`)) {
		return unquote(s)
	}
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return parseFlowSequence(s)
	}
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		return parseFlowMapping(s)
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func parseFlowSequence(s string) []any {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '[' || s[len(s)-1] != ']' {
		return nil
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return []any{}
	}
	parts := splitFlowItems(inner)
	res := make([]any, 0, len(parts))
	for _, p := range parts {
		res = append(res, parseYAMLValue(strings.TrimSpace(p)))
	}
	return res
}

func parseFlowMapping(s string) map[string]any {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return nil
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		return map[string]any{}
	}
	parts := splitFlowItems(inner)
	res := make(map[string]any, len(parts))
	for _, p := range parts {
		colon := strings.IndexByte(p, ':')
		if colon > 0 {
			k := strings.TrimSpace(p[:colon])
			v := strings.TrimSpace(p[colon+1:])
			res[unquote(k)] = parseYAMLValue(v)
		}
	}
	return res
}

func splitFlowItems(s string) []string {
	var parts []string
	bracketDepth := 0
	braceDepth := 0
	inQuotes := false
	quoteChar := byte(0)
	start := 0

	for i := 0; i < len(s); i++ {
		b := s[i]
		if inQuotes {
			if b == quoteChar && (i == 0 || s[i-1] != '\\') {
				inQuotes = false
			}
			continue
		}
		if b == '"' || b == '\'' {
			inQuotes = true
			quoteChar = b
			continue
		}
		if b == '[' {
			bracketDepth++
		} else if b == ']' {
			bracketDepth--
		} else if b == '{' {
			braceDepth++
		} else if b == '}' {
			braceDepth--
		} else if b == ',' && bracketDepth == 0 && braceDepth == 0 {
			parts = append(parts, strings.TrimSpace(s[start:i]))
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, strings.TrimSpace(s[start:]))
	}
	return parts
}

func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
