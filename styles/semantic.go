package styles

import "strings"

// SemanticRule holds pre-designed semantic component definitions for the Code Garden aesthetic.
type SemanticComponent struct {
	Name    string
	Trigger string // class trigger, e.g. "btn"
	CSS     string
}

// SemanticComponents returns all semantic components in KrewireCSS.
var SemanticComponents = []SemanticComponent{
	{
		Name:    "Badge",
		Trigger: "badge",
		CSS: `.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  border-radius: 9999px;
  padding: 0.2rem 0.65rem;
  font-family: var(--ui-font-sans);
  font-size: var(--ui-text-xs);
  line-height: var(--ui-leading-snug);
  font-weight: var(--ui-font-semibold);
  letter-spacing: var(--ui-tracking-wide);
  border: 1px solid var(--border);
  box-shadow: none;
  white-space: nowrap;
}
.badge-primary {
  background-color: var(--primary-subtle);
  color: var(--primary);
  border-color: rgba(57, 211, 83, 0.3);
}
.badge-secondary {
  background-color: var(--secondary-subtle);
  color: var(--secondary);
  border-color: rgba(0, 209, 193, 0.3);
}
.badge-accent {
  background-color: var(--accent-subtle);
  color: var(--accent);
  border-color: rgba(255, 59, 46, 0.3);
}
.badge-muted {
  background-color: var(--bg-2);
  color: var(--fg-muted);
}
.badge-green {
  background-color: rgba(57, 211, 83, 0.14);
  border-color: rgba(57, 211, 83, 0.3);
  color: var(--primary, #39D353);
}
.badge-cyan {
  background-color: rgba(0, 209, 193, 0.14);
  border-color: rgba(0, 209, 193, 0.3);
  color: var(--secondary, #00D1C1);
}
.badge-yellow {
  background-color: rgba(245, 158, 11, 0.12);
  border-color: rgba(245, 158, 11, 0.3);
  color: var(--accent, #f59e0b);
}
.badge-red {
  background-color: rgba(255, 59, 46, 0.14);
  border-color: rgba(255, 59, 46, 0.3);
  color: var(--error, #FF3B2E);
}
.badge-purple {
  background-color: rgba(139, 92, 246, 0.12);
  border-color: rgba(139, 92, 246, 0.3);
  color: #8b5cf6;
}
.badge-gray {
  background-color: var(--bg-2);
  border-color: var(--border);
  color: var(--fg-muted);
}
.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  display: inline-block;
}`,
	},
	{
		Name:    "Card",
		Trigger: "card",
		CSS: `.card {
  display: flex;
  width: 100%;
  min-width: 0;
  flex-direction: column;
  border-radius: var(--ui-radius-lg, 12px);
  background-color: var(--surface);
  padding: 1.25rem;
  border: 1px solid var(--border);
  box-shadow: var(--ui-shadow-sm);
  transition: transform var(--ui-duration-fast) var(--ui-ease-standard), box-shadow var(--ui-duration-fast) var(--ui-ease-standard), border-color var(--ui-duration-fast) var(--ui-ease-standard);
  box-sizing: border-box;
}
.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--ui-shadow-md);
  border-color: var(--border-focus);
}
.card-clickable {
  text-decoration: none;
  color: inherit;
  cursor: pointer;
}
.card-featured {
  background: var(--bg-2);
  border-color: var(--primary);
  box-shadow: var(--ui-shadow);
}
.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}
.card-titles {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}
.card-title {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 700;
  letter-spacing: -0.01em;
  color: inherit;
}
.card-subtitle {
  margin: 0;
  font-size: 0.85rem;
  color: var(--fg-muted);
  font-weight: 400;
}
.card-badge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.2rem 0.55rem;
}
.card-header {
  font-weight: 700;
  font-size: 1.15rem;
  margin-bottom: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--border);
  padding-bottom: 0.5rem;
}`,
	},
	{
		Name:    "Button",
		Trigger: "btn",
		CSS: `.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-family: inherit;
  font-weight: 600;
  border-radius: var(--radius, 8px);
  border: 1px solid var(--border);
  box-shadow: var(--ui-shadow-sm);
  cursor: pointer;
  text-decoration: none;
  transition: transform 0.15s ease, box-shadow 0.15s ease, background-color 0.15s ease, border-color 0.15s ease;
  line-height: 1.25;
  white-space: nowrap;
  box-sizing: border-box;
}
.btn:hover:not(:disabled):not([aria-disabled="true"]) {
  transform: translateY(-1px);
  box-shadow: var(--ui-shadow);
}
.btn:active:not(:disabled):not([aria-disabled="true"]) {
  transform: translateY(0px);
  box-shadow: var(--ui-shadow-sm);
}
.btn:disabled, .btn[aria-disabled="true"] {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}
.btn-xs { padding: 0.25rem 0.5rem; font-size: 0.75rem; border-radius: var(--radius-sm, 6px); }
.btn-sm { padding: 0.375rem 0.75rem; font-size: 0.8125rem; border-radius: var(--radius-sm, 6px); }
.btn-md { padding: 0.5rem 1rem; font-size: 0.875rem; border-radius: var(--radius, 8px); }
.btn-lg { padding: 0.625rem 1.25rem; font-size: 1rem; border-radius: var(--radius-lg, 12px); }

.btn-primary, button.btn-primary, a.btn-primary {
  background: var(--primary, #39D353);
  color: var(--primary-content, #0B1F3B) !important;
  border-color: var(--primary, #39D353);
}
.btn-primary:hover:not(:disabled):not([aria-disabled="true"]),
button.btn-primary:hover:not(:disabled):not([aria-disabled="true"]),
a.btn-primary:hover:not(:disabled):not([aria-disabled="true"]) {
  background: var(--primary-hover, #2ecc50);
  border-color: var(--primary-hover, #2ecc50);
  color: var(--primary-content, #0B1F3B) !important;
}

.btn-secondary, button.btn-secondary, a.btn-secondary {
  background: var(--secondary, #00D1C1);
  color: var(--secondary-content, #0B1F3B) !important;
  border-color: var(--secondary, #00D1C1);
}
.btn-secondary:hover:not(:disabled):not([aria-disabled="true"]),
button.btn-secondary:hover:not(:disabled):not([aria-disabled="true"]),
a.btn-secondary:hover:not(:disabled):not([aria-disabled="true"]) {
  background: var(--secondary-hover, #00b8aa);
  border-color: var(--secondary-hover, #00b8aa);
  color: var(--secondary-content, #0B1F3B) !important;
}

.btn-accent, .btn-danger, .btn-error,
button.btn-accent, button.btn-danger, button.btn-error,
a.btn-accent, a.btn-danger, a.btn-error {
  background: var(--accent, #FF3B2E);
  color: #ffffff !important;
  border-color: var(--accent, #FF3B2E);
}
.btn-accent:hover:not(:disabled):not([aria-disabled="true"]),
button.btn-accent:hover:not(:disabled):not([aria-disabled="true"]),
a.btn-accent:hover:not(:disabled):not([aria-disabled="true"]) {
  background: var(--accent-hover, #e0281b);
  border-color: var(--accent-hover, #e0281b);
  color: #ffffff !important;
}

.btn-outline, button.btn-outline, a.btn-outline {
  background: transparent;
  color: var(--fg, #0B1F3B) !important;
  border-color: var(--border, #e2ded4);
}
.btn-outline:hover:not(:disabled):not([aria-disabled="true"]),
button.btn-outline:hover:not(:disabled):not([aria-disabled="true"]),
a.btn-outline:hover:not(:disabled):not([aria-disabled="true"]) {
  background: var(--bg-2, #F2EDE4);
  color: var(--fg, #0B1F3B) !important;
}

.btn-ghost, button.btn-ghost, a.btn-ghost {
  background: transparent;
  border-color: transparent;
  box-shadow: none;
  color: var(--fg, #0B1F3B) !important;
}
.btn-ghost:hover:not(:disabled):not([aria-disabled="true"]),
button.btn-ghost:hover:not(:disabled):not([aria-disabled="true"]),
a.btn-ghost:hover:not(:disabled):not([aria-disabled="true"]) {
  background: var(--bg-2, #f5f5f4);
  color: var(--fg, #1c1917) !important;
  transform: none;
  box-shadow: none;
}

.btn-logout {
  background: var(--surface, #ffffff);
  border: 1px solid var(--border, #e7e5e4);
  font-weight: 600;
  padding: 0.35rem 0.75rem;
  border-radius: var(--radius-sm, 6px);
  cursor: pointer;
  box-shadow: var(--ui-shadow-sm);
  font-family: inherit;
  color: var(--fg, #1c1917);
  transition: all 0.15s ease;
}
.btn-logout:hover {
  background: var(--error, #FF3B2E);
  border-color: var(--error, #FF3B2E);
  color: #ffffff;
}

.btn-sso {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.625rem;
  width: 100%;
  padding: 0.75rem 1.25rem;
  background-color: var(--primary, #39D353);
  color: var(--primary-content, #0B1F3B);
  font-family: inherit;
  font-size: 0.95rem;
  font-weight: 600;
  text-decoration: none;
  border: 1px solid var(--primary, #39D353);
  border-radius: var(--radius, 8px);
  box-shadow: var(--ui-shadow-sm);
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease, background-color 0.15s ease;
  box-sizing: border-box;
}
.btn-sso:hover {
  background-color: var(--primary-hover, #2ecc50);
  transform: translateY(-1px);
  box-shadow: var(--ui-shadow);
  color: var(--primary-content, #0B1F3B);
  text-decoration: none;
}
.btn-sso:active {
  transform: translateY(0px);
  box-shadow: var(--ui-shadow-sm);
}`,
	},
	{
		Name:    "Alert",
		Trigger: "alert",
		CSS: `.alert, .alert-banner {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.85rem 1.15rem;
  border-radius: var(--radius, 8px);
  border: 1px solid var(--border, #e7e5e4);
  box-shadow: var(--ui-shadow-sm);
  margin: 1rem 0;
  font-size: 0.875rem;
  line-height: 1.5;
  box-sizing: border-box;
}
.alert-icon {
  width: 20px;
  height: 20px;
  font-size: 1.1rem;
  font-weight: 700;
  line-height: 1.2;
  flex-shrink: 0;
}
.alert-content {
  flex: 1;
  min-width: 0;
}
.alert-title {
  display: block;
  font-size: 0.9rem;
  font-weight: 600;
  margin-bottom: 0.2rem;
}
.alert-message {
  font-size: 0.85rem;
}
.alert-close {
  background: transparent;
  border: 1px solid var(--border, #e7e5e4);
  border-radius: 6px;
  cursor: pointer;
  font-size: 1.1rem;
  line-height: 1;
  padding: 0.15rem 0.4rem;
  color: var(--fg, #1c1917);
}
.alert-info {
  background: rgba(13, 148, 136, 0.08);
  border-color: rgba(13, 148, 136, 0.2);
  color: var(--fg, #1c1917);
}
.alert-success {
  background: rgba(16, 185, 129, 0.08);
  border-color: rgba(16, 185, 129, 0.2);
  color: var(--fg, #1c1917);
}
.alert-warning {
  background: rgba(245, 158, 11, 0.08);
  border-color: rgba(245, 158, 11, 0.2);
  color: var(--fg, #1c1917);
}
.alert-error, .alert-danger {
  background-color: rgba(239, 68, 68, 0.08);
  border-color: rgba(239, 68, 68, 0.25);
  color: #991b1b;
  box-shadow: var(--ui-shadow-sm);
}
.alert-error .alert-icon, .alert-danger .alert-icon {
  color: #ef4444;
}`,
	},
	{
		Name:    "Table",
		Trigger: "table",
		CSS: `.table-wrapper, .kiw-table-wrapper {
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border-radius: var(--radius-lg, 12px);
  border: 1px solid var(--border, #e7e5e4);
  box-sizing: border-box;
  margin-bottom: 1rem;
}
table.brutalist-table, table.kiw-table, table.table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
  font-size: 0.875rem;
  color: var(--fg, #1c1917);
  background: var(--surface, #ffffff);
}
table.brutalist-table th, table.kiw-table th, table.table th {
  background: var(--bg-2, #f5f5f4);
  color: var(--fg-muted, #78716c);
  border-bottom: 1px solid var(--border, #e7e5e4);
  padding: 0.75rem 1rem;
  font-weight: 600;
  font-size: 0.75rem;
  letter-spacing: 0.03em;
  white-space: nowrap;
}
table.brutalist-table td, table.kiw-table td, table.table td {
  border-bottom: 1px solid var(--border, #e7e5e4);
  padding: 0.75rem 1rem;
  vertical-align: middle;
}
table.brutalist-table tr:last-child td,
table.kiw-table tr:last-child td,
table.table tr:last-child td {
  border-bottom: none;
}
table.brutalist-table tr:hover td,
table.kiw-table tr:hover td,
table.table tr:hover td {
  background-color: rgba(16, 185, 129, 0.04);
}`,
	},
	{
		Name:    "Form",
		Trigger: "form",
		CSS: `.form-group, .field-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  margin-bottom: 1.25rem;
  box-sizing: border-box;
}
label.form-label, .field-label {
  display: block;
  font-weight: 600;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
  letter-spacing: normal;
  color: var(--fg, #1c1917);
}
.field-required {
  color: var(--error, #ef4444);
  font-weight: 700;
  margin-left: 0.2rem;
}
input.form-input, .field-input,
select.form-select, .field-select,
textarea.form-textarea, .field-textarea {
  width: 100%;
  padding: 0.55rem 0.875rem;
  font-size: 0.875rem;
  font-family: inherit;
  border-radius: var(--radius, 8px);
  border: 1px solid var(--border, #e2ded4);
  background: var(--surface, #ffffff);
  color: var(--fg, #0B1F3B);
  box-shadow: var(--ui-shadow-sm);
  outline: none;
  box-sizing: border-box;
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}
input.form-input:focus, .field-input:focus,
select.form-select:focus, textarea.form-textarea:focus {
  border-color: var(--primary, #39D353);
  box-shadow: 0 0 0 2px rgba(57, 211, 83, 0.2);
}
.field-input.has-error, input.form-input.has-error {
  border-color: var(--error, #FF3B2E);
}
.form-hint, .field-help {
  font-size: 0.78rem;
  color: var(--fg-muted, #57677d);
  margin-top: 0.25rem;
}
.form-error, .field-error {
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--error, #FF3B2E);
  margin-top: 0.25rem;
}`,
	},
	{
		Name:    "AppShell",
		Trigger: "top-bar",
		CSS: `header.top-bar, .top-bar {
  background-color: var(--surface, #ffffff);
  border: none;
  border-bottom: 1px solid var(--border, #e7e5e4);
  box-shadow: var(--ui-shadow-sm);
  padding: 0.875rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: var(--ui-z-sticky, 100);
  box-sizing: border-box;
}
.brand-icon-box {
  display: inline-grid;
  place-items: center;
  width: 32px;
  height: 32px;
  background: var(--primary, #39D353);
  border: 1px solid var(--primary, #39D353);
  border-radius: var(--radius-sm, 6px);
  box-shadow: var(--ui-shadow-sm);
  color: var(--primary-content, #0B1F3B);
  flex-shrink: 0;
}
.brand-badge {
  background: rgba(0, 209, 193, 0.1);
  border: 1px solid rgba(0, 209, 193, 0.3);
  border-radius: 9999px;
  padding: 0.15rem 0.55rem;
  font-size: 0.72rem;
  font-family: var(--ui-font-mono);
  font-weight: 600;
  box-shadow: none;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--secondary, #0d9488);
}
.metric-card {
  background: var(--surface, #ffffff);
  border: 1px solid var(--border, #e7e5e4);
  border-radius: var(--radius-lg, 12px);
  box-shadow: var(--ui-shadow-sm);
  padding: 1.25rem;
  position: relative;
  transition: transform 0.15s ease, box-shadow 0.15s ease, border-color 0.15s ease;
  box-sizing: border-box;
}
.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--ui-shadow-md);
  border-color: var(--border-focus);
}`,
	},
	{
		Name:    "Navbar",
		Trigger: "navbar",
		CSS: `.navbar,
header.navbar {
  border: none !important;
  border-top: none !important;
  border-left: none !important;
  border-right: none !important;
  border-bottom: 1px solid var(--border) !important;
  box-shadow: var(--ui-shadow-sm) !important;
  box-sizing: border-box;
}
:root[data-theme="dark"] .navbar,
html[data-theme="dark"] .navbar,
.dark .navbar {
  border-bottom: 1px solid var(--border) !important;
}
.navbar .container-site {
  padding-top: 0.875rem !important;
  padding-bottom: 0.875rem !important;
}
.nav-menu {
  display: none;
}
@media (min-width: 768px) {
  .nav-menu,
  div.nav-menu,
  nav.nav-menu,
  .navbar .nav-menu {
    display: flex !important;
    align-items: center;
  }
}`,
	},
}

// RenderSemanticCSS generates the CSS rules for semantic components whose triggers are matched.
func RenderSemanticCSS(usedTriggers map[string]bool) string {
	var b strings.Builder
	for _, comp := range SemanticComponents {
		if usedTriggers[comp.Trigger] || usedTriggers[strings.ToLower(comp.Name)] {
			b.WriteString(comp.CSS)
			b.WriteString("\n\n")
		}
	}
	return b.String()
}
