package styles

import (
	"fmt"
	"strings"
)

type Rule struct {
	Selector    string
	Declaration string
	Media       string // e.g. "min-width: 640px" for sm
	Prefix      string // e.g. "sm:", "hover:", "dark:"
}

// Spacing scale mapping
var spacingScale = map[string]string{
	"0":    "0px",
	"0.5":  "0.125rem",
	"1":    "0.25rem",
	"1.5":  "0.375rem",
	"2":    "0.5rem",
	"2.5":  "0.625rem",
	"3":    "0.75rem",
	"3.5":  "0.875rem",
	"4":    "1rem",
	"5":    "1.25rem",
	"6":    "1.5rem",
	"8":    "2rem",
	"10":   "2.5rem",
	"12":   "3rem",
	"16":   "4rem",
	"20":   "5rem",
	"24":   "6rem",
	"auto": "auto",
}

// Breakpoint scale
var breakpoints = map[string]string{
	"sm":  "640px",
	"md":  "768px",
	"lg":  "1024px",
	"xl":  "1280px",
	"2xl": "1536px",
}

// ResolveUtility parses a class name and generates a CSS Rule if it is a valid utility.
func ResolveUtility(rawClass string) (*Rule, bool) {
	class := rawClass
	media := ""
	pseudo := ""
	isDark := false

	// Check prefixes (e.g. sm:hover:bg-primary or md:flex)
	for {
		if strings.HasPrefix(class, "sm:") {
			media = "min-width: 640px"
			class = strings.TrimPrefix(class, "sm:")
		} else if strings.HasPrefix(class, "md:") {
			media = "min-width: 768px"
			class = strings.TrimPrefix(class, "md:")
		} else if strings.HasPrefix(class, "lg:") {
			media = "min-width: 1024px"
			class = strings.TrimPrefix(class, "lg:")
		} else if strings.HasPrefix(class, "xl:") {
			media = "min-width: 1280px"
			class = strings.TrimPrefix(class, "xl:")
		} else if strings.HasPrefix(class, "2xl:") {
			media = "min-width: 1536px"
			class = strings.TrimPrefix(class, "2xl:")
		} else if strings.HasPrefix(class, "hover:") {
			pseudo = ":hover"
			class = strings.TrimPrefix(class, "hover:")
		} else if strings.HasPrefix(class, "focus:") {
			pseudo = ":focus"
			class = strings.TrimPrefix(class, "focus:")
		} else if strings.HasPrefix(class, "active:") {
			pseudo = ":active"
			class = strings.TrimPrefix(class, "active:")
		} else if strings.HasPrefix(class, "disabled:") {
			pseudo = ":disabled"
			class = strings.TrimPrefix(class, "disabled:")
		} else if strings.HasPrefix(class, "dark:") {
			isDark = true
			class = strings.TrimPrefix(class, "dark:")
		} else {
			break
		}
	}

	decl := resolveBaseUtility(class)
	if decl == "" {
		return nil, false
	}

	escapedSelector := escapeSelector(rawClass)
	if pseudo != "" {
		escapedSelector = escapedSelector + pseudo
	}
	if isDark {
		escapedSelector = fmt.Sprintf("[data-theme=\"dark\"] .%s, .dark .%s", escapedSelector, escapedSelector)
	} else {
		escapedSelector = "." + escapedSelector
	}

	return &Rule{
		Selector:    escapedSelector,
		Declaration: decl,
		Media:       media,
	}, true
}

func escapeSelector(cls string) string {
	res := strings.ReplaceAll(cls, ":", `\:`)
	res = strings.ReplaceAll(res, "/", `\/`)
	res = strings.ReplaceAll(res, ".", `\.`)
	return res
}

func resolveBaseUtility(cls string) string {
	// Layout
	switch cls {
	case "flex":
		return "display: flex;"
	case "inline-flex":
		return "display: inline-flex;"
	case "grid":
		return "display: grid;"
	case "block":
		return "display: block;"
	case "inline-block":
		return "display: inline-block;"
	case "inline":
		return "display: inline;"
	case "hidden":
		return "display: none;"
	case "flex-col":
		return "flex-direction: column;"
	case "flex-col-reverse":
		return "flex-direction: column-reverse;"
	case "flex-row":
		return "flex-direction: row;"
	case "flex-row-reverse":
		return "flex-direction: row-reverse;"
	case "flex-wrap":
		return "flex-wrap: wrap;"
	case "flex-nowrap":
		return "flex-wrap: nowrap;"
	case "items-center":
		return "align-items: center;"
	case "items-start":
		return "align-items: flex-start;"
	case "items-end":
		return "align-items: flex-end;"
	case "items-baseline":
		return "align-items: baseline;"
	case "items-stretch":
		return "align-items: stretch;"
	case "justify-center":
		return "justify-content: center;"
	case "justify-start":
		return "justify-content: flex-start;"
	case "justify-end":
		return "justify-content: flex-end;"
	case "justify-between":
		return "justify-content: space-between;"
	case "justify-around":
		return "justify-content: space-around;"
	case "justify-evenly":
		return "justify-content: space-evenly;"
	case "flex-1":
		return "flex: 1 1 0%;"
	case "flex-auto":
		return "flex: 1 1 auto;"
	case "flex-none":
		return "flex: none;"
	case "shrink-0":
		return "flex-shrink: 0;"
	case "grow":
		return "flex-grow: 1;"
	case "relative":
		return "position: relative;"
	case "absolute":
		return "position: absolute;"
	case "sticky":
		return "position: sticky;"
	case "fixed":
		return "position: fixed;"
	case "top-0":
		return "top: 0px;"
	case "bottom-0":
		return "bottom: 0px;"
	case "left-0":
		return "left: 0px;"
	case "right-0":
		return "right: 0px;"
	case "inset-0":
		return "inset: 0px;"
	case "z-10":
		return "z-index: 10;"
	case "z-20":
		return "z-index: 20;"
	case "z-50":
		return "z-index: 50;"
	}

	// Grid
	if strings.HasPrefix(cls, "grid-cols-") {
		cols := strings.TrimPrefix(cls, "grid-cols-")
		return fmt.Sprintf("grid-template-columns: repeat(%s, minmax(0, 1fr));", cols)
	}
	if strings.HasPrefix(cls, "grid-") {
		cols := strings.TrimPrefix(cls, "grid-")
		return fmt.Sprintf("grid-template-columns: repeat(%s, minmax(0, 1fr));", cols)
	}
	if strings.HasPrefix(cls, "col-span-") {
		span := strings.TrimPrefix(cls, "col-span-")
		return fmt.Sprintf("grid-column: span %s / span %s;", span, span)
	}

	// Sizing
	switch cls {
	case "w-full":
		return "width: 100%;"
	case "w-screen":
		return "width: 100vw;"
	case "w-auto":
		return "width: auto;"
	case "w-fit":
		return "width: fit-content;"
	case "h-full":
		return "height: 100%;"
	case "h-screen":
		return "height: 100vh;"
	case "h-auto":
		return "height: auto;"
	case "min-h-screen":
		return "min-height: 100vh;"
	case "min-w-0":
		return "min-width: 0px;"
	case "max-w-xs":
		return "max-width: 20rem;"
	case "max-w-sm":
		return "max-width: 24rem;"
	case "max-w-md":
		return "max-width: 28rem;"
	case "max-w-lg":
		return "max-width: 32rem;"
	case "max-w-xl":
		return "max-width: 36rem;"
	case "max-w-2xl":
		return "max-width: 42rem;"
	case "max-w-site":
		return "max-width: 108rem;"
	case "container":
		return "width: 100%; max-width: 1240px; margin-left: auto; margin-right: auto; padding-left: 1rem; padding-right: 1rem;"
	case "container-fluid":
		return "width: 100%; padding-left: 1rem; padding-right: 1rem; margin-left: auto; margin-right: auto;"
	case "section-py":
		return "padding-top: 3rem; padding-bottom: 3rem;"
	}

	// Spacing: p, px, py, pt, pr, pb, pl, m, mx, my, mt, mr, mb, ml, gap, gap-x, gap-y
	parts := strings.SplitN(cls, "-", 2)
	if len(parts) == 2 {
		prefix := parts[0]
		valKey := parts[1]
		valKeyNorm := strings.ReplaceAll(valKey, "_", ".")
		if val, ok := spacingScale[valKeyNorm]; ok {
			switch prefix {
			case "p":
				return fmt.Sprintf("padding: %s;", val)
			case "px":
				return fmt.Sprintf("padding-left: %s; padding-right: %s;", val, val)
			case "py":
				return fmt.Sprintf("padding-top: %s; padding-bottom: %s;", val, val)
			case "pt":
				return fmt.Sprintf("padding-top: %s;", val)
			case "pr":
				return fmt.Sprintf("padding-right: %s;", val)
			case "pb":
				return fmt.Sprintf("padding-bottom: %s;", val)
			case "pl":
				return fmt.Sprintf("padding-left: %s;", val)
			case "m":
				return fmt.Sprintf("margin: %s;", val)
			case "mx":
				return fmt.Sprintf("margin-left: %s; margin-right: %s;", val, val)
			case "my":
				return fmt.Sprintf("margin-top: %s; margin-bottom: %s;", val, val)
			case "mt":
				return fmt.Sprintf("margin-top: %s;", val)
			case "mr":
				return fmt.Sprintf("margin-right: %s;", val)
			case "mb":
				return fmt.Sprintf("margin-bottom: %s;", val)
			case "ml":
				return fmt.Sprintf("margin-left: %s;", val)
			case "gap":
				return fmt.Sprintf("gap: %s;", val)
			case "gap-x":
				return fmt.Sprintf("column-gap: %s;", val)
			case "gap-y":
				return fmt.Sprintf("row-gap: %s;", val)
			}
		}
	}

	// Typography
	switch cls {
	case "text-2xs":
		return "font-size: 0.625rem;"
	case "text-xs":
		return "font-size: 0.75rem;"
	case "text-sm":
		return "font-size: 0.875rem;"
	case "text-base":
		return "font-size: 1rem;"
	case "text-md":
		return "font-size: 1rem;"
	case "text-lg":
		return "font-size: 1.125rem;"
	case "text-xl":
		return "font-size: 1.25rem;"
	case "text-2xl":
		return "font-size: 1.5rem;"
	case "text-3xl":
		return "font-size: 1.875rem;"
	case "text-4xl":
		return "font-size: 2.25rem;"
	case "text-5xl":
		return "font-size: 3rem;"
	case "font-normal":
		return "font-weight: 400;"
	case "font-medium":
		return "font-weight: 500;"
	case "font-semibold":
		return "font-weight: 600;"
	case "font-bold":
		return "font-weight: 700;"
	case "font-extrabold":
		return "font-weight: 800;"
	case "font-black":
		return "font-weight: 900;"
	case "font-sans":
		return "font-family: var(--ui-font-sans);"
	case "font-mono":
		return "font-family: var(--ui-font-mono);"
	case "font-display":
		return "font-family: var(--ui-font-display);"
	case "text-left":
		return "text-align: left;"
	case "text-center":
		return "text-align: center;"
	case "text-right":
		return "text-align: right;"
	case "uppercase":
		return "text-transform: uppercase;"
	case "lowercase":
		return "text-transform: lowercase;"
	case "capitalize":
		return "text-transform: capitalize;"
	case "truncate":
		return "overflow: hidden; text-overflow: ellipsis; white-space: nowrap;"
	case "whitespace-nowrap":
		return "white-space: nowrap;"
	case "leading-none":
		return "line-height: 1;"
	case "leading-tight":
		return "line-height: 1.25;"
	case "leading-normal":
		return "line-height: 1.5;"
	case "leading-relaxed":
		return "line-height: 1.625;"
	case "tracking-tight":
		return "letter-spacing: -0.025em;"
	case "tracking-normal":
		return "letter-spacing: 0em;"
	case "tracking-wide":
		return "letter-spacing: 0.025em;"
	case "tracking-wider":
		return "letter-spacing: 0.05em;"
	case "tracking-widest":
		return "letter-spacing: 0.1em;"
	}

	// Colors
	switch cls {
	case "text-primary":
		return "color: var(--primary);"
	case "text-secondary":
		return "color: var(--secondary);"
	case "text-accent":
		return "color: var(--accent);"
	case "text-fg":
		return "color: var(--fg);"
	case "text-muted":
		return "color: var(--fg-muted);"
	case "text-white":
		return "color: #ffffff;"
	case "text-error":
		return "color: var(--error);"
	case "bg-primary":
		return "background-color: var(--primary);"
	case "bg-secondary":
		return "background-color: var(--secondary);"
	case "bg-accent":
		return "background-color: var(--accent);"
	case "bg-surface":
		return "background-color: var(--surface);"
	case "bg-bg":
		return "background-color: var(--bg);"
	case "bg-muted":
		return "background-color: var(--bg-2);"
	case "bg-white":
		return "background-color: #ffffff;"
	case "bg-transparent":
		return "background-color: transparent;"
	case "bg-error":
		return "background-color: var(--error);"
	}

	// Borders & Radius
	switch cls {
	case "border":
		return "border: 1px solid var(--border);"
	case "border-0":
		return "border-width: 0px;"
	case "border-2":
		return "border-width: 2px;"
	case "border-t":
		return "border-top: 1px solid var(--border);"
	case "border-b":
		return "border-bottom: 1px solid var(--border);"
	case "border-l":
		return "border-left: 1px solid var(--border);"
	case "border-r":
		return "border-right: 1px solid var(--border);"
	case "border-border":
		return "border-color: var(--border);"
	case "border-primary":
		return "border-color: var(--primary);"
	case "border-secondary":
		return "border-color: var(--secondary);"
	case "border-transparent":
		return "border-color: transparent;"
	case "border-pop":
		return "border: 1px solid var(--border);"
	case "border-pop-thin":
		return "border: 1px solid var(--border-subtle, var(--border));"
	case "rounded-none":
		return "border-radius: 0px;"
	case "rounded-xs":
		return "border-radius: var(--ui-radius-xs, 4px);"
	case "rounded-sm":
		return "border-radius: var(--radius-sm, 6px);"
	case "rounded":
		return "border-radius: var(--radius, 8px);"
	case "rounded-md":
		return "border-radius: var(--radius-md, 10px);"
	case "rounded-lg":
		return "border-radius: var(--radius-lg, 12px);"
	case "rounded-xl":
		return "border-radius: var(--radius-xl, 16px);"
	case "rounded-2xl":
		return "border-radius: 1.25rem;"
	case "rounded-full":
		return "border-radius: 9999px;"
	}

	// Shadows
	switch cls {
	case "shadow-none":
		return "box-shadow: none;"
	case "shadow-sm":
		return "box-shadow: var(--ui-shadow-sm);"
	case "shadow":
		return "box-shadow: var(--ui-shadow);"
	case "shadow-md":
		return "box-shadow: var(--ui-shadow-md);"
	case "shadow-lg":
		return "box-shadow: var(--ui-shadow-lg);"
	case "shadow-xl":
		return "box-shadow: var(--ui-shadow-xl);"
	case "shadow-pop":
		return "box-shadow: var(--shadow-sm, 0 1px 2px 0 rgba(0,0,0,0.05));"
	case "shadow-pop-sm":
		return "box-shadow: var(--shadow-sm, 0 1px 2px 0 rgba(0,0,0,0.05));"
	case "shadow-pop-lg":
		return "box-shadow: var(--shadow-md, 0 4px 6px -1px rgba(0,0,0,0.08));"
	}

	// Transitions & Interactivity
	switch cls {
	case "transition-all":
		return "transition-property: all; transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1); transition-duration: 150ms;"
	case "transition-colors":
		return "transition-property: color, background-color, border-color; transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1); transition-duration: 150ms;"
	case "transition-transform":
		return "transition-property: transform; transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1); transition-duration: 150ms;"
	case "duration-150":
		return "transition-duration: 150ms;"
	case "duration-200":
		return "transition-duration: 200ms;"
	case "duration-300":
		return "transition-duration: 300ms;"
	case "cursor-pointer":
		return "cursor: pointer;"
	case "cursor-not-allowed":
		return "cursor: not-allowed;"
	case "opacity-50":
		return "opacity: 0.5;"
	case "opacity-100":
		return "opacity: 1;"
	}

	return ""
}
