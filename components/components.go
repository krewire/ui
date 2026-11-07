package components

import (
	"embed"
	"fmt"
	"sort"
	"strings"
)

// FS embeds all built-in .kiw reusable and atomic UI components.
//
//go:embed *.kiw
var FS embed.FS

// ComponentMeta represents the metadata and category for a .kiw UI component.
type ComponentMeta struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// Registry maps each component name to its metadata and classification.
var Registry = map[string]ComponentMeta{
	// Layout
	"Box":       {Name: "Box", Category: "Layout", Description: "Generic layout container with customizable padding, margin, and borders"},
	"Container": {Name: "Container", Category: "Layout", Description: "Responsive max-width page container"},
	"Flex":      {Name: "Flex", Category: "Layout", Description: "Flexible box layout container with direction, wrap, and alignment"},
	"Grid":      {Name: "Grid", Category: "Layout", Description: "Responsive CSS grid container with customizable columns and gaps"},
	"HStack":    {Name: "HStack", Category: "Layout", Description: "Horizontal flex stack with configurable spacing and alignment"},
	"Section":   {Name: "Section", Category: "Layout", Description: "Semantic section wrapper with vertical rhythm"},
	"Stack":     {Name: "Stack", Category: "Layout", Description: "Directional stack wrapper with uniform spacing"},
	"VStack":    {Name: "VStack", Category: "Layout", Description: "Vertical flex stack with configurable spacing and alignment"},

	// Atoms & Primitives
	"Avatar":   {Name: "Avatar", Category: "Atoms", Description: "Image or initials avatar with size variants and status indicators"},
	"Badge":    {Name: "Badge", Category: "Atoms", Description: "Compact label for status, counts, or categorization"},
	"Brand":    {Name: "Brand", Category: "Atoms", Description: "App brand logo and name display"},
	"Button":   {Name: "Button", Category: "Atoms", Description: "Interactive button with primary, secondary, ghost, and danger variants"},
	"Code":     {Name: "Code", Category: "Atoms", Description: "Inline code badge with monospace typography"},
	"Divider":  {Name: "Divider", Category: "Atoms", Description: "Horizontal or vertical separator line"},
	"Heading":  {Name: "Heading", Category: "Atoms", Description: "Semantic heading element (h1-h6) with consistent typography tokens"},
	"Kbd":      {Name: "Kbd", Category: "Atoms", Description: "Keyboard shortcut keycap indicator"},
	"Link":     {Name: "Link", Category: "Atoms", Description: "Accessible hyperlink with hover states and external indicator"},
	"Progress": {Name: "Progress", Category: "Atoms", Description: "Determinate or indeterminate progress bar indicator"},
	"Spinner":  {Name: "Spinner", Category: "Atoms", Description: "Loading spinner indicator with speed and color variants"},
	"Text":     {Name: "Text", Category: "Atoms", Description: "Versatile typography text atom with muted, lead, and size props"},

	// Forms & Inputs
	"Checkbox": {Name: "Checkbox", Category: "Forms", Description: "Accessible checkbox input with label and description support"},
	"Form":     {Name: "Form", Category: "Forms", Description: "Form layout wrapper with accessible grouping and validation slots"},
	"Input":    {Name: "Input", Category: "Forms", Description: "Text input with prefix/suffix icons, helper text, and error states"},
	"Radio":    {Name: "Radio", Category: "Forms", Description: "Radio selection control for mutually exclusive options"},
	"Select":   {Name: "Select", Category: "Forms", Description: "Dropdown selection menu with custom styled options"},
	"Textarea": {Name: "Textarea", Category: "Forms", Description: "Multi-line text input with autoresize and character counter"},
	"Toggle":   {Name: "Toggle", Category: "Forms", Description: "Binary switch toggle for boolean settings"},

	// Navigation
	"Breadcrumb":     {Name: "Breadcrumb", Category: "Navigation", Description: "Hierarchical page breadcrumb navigation trail"},
	"BreadcrumbItem": {Name: "BreadcrumbItem", Category: "Navigation", Description: "Individual crumb segment within a Breadcrumb navigation"},
	"Hamburger":      {Name: "Hamburger", Category: "Navigation", Description: "Animated mobile navigation toggle button"},
	"LangSwitch":     {Name: "LangSwitch", Category: "Navigation", Description: "Multi-language locale switcher dropdown"},
	"Navbar":         {Name: "Navbar", Category: "Navigation", Description: "Top navigation bar with brand, menu links, slots, and responsive drawer"},
	"NavMenu":        {Name: "NavMenu", Category: "Navigation", Description: "Composable navigation item list with slot for active page tracking"},
	"Sidebar":        {Name: "Sidebar", Category: "Navigation", Description: "Full-height vertical collapsible side navigation container"},
	"SideMenu":       {Name: "SideMenu", Category: "Navigation", Description: "Structured sidebar links with category groupings and badges"},
	"Tabs":           {Name: "Tabs", Category: "Navigation", Description: "Tabbed content switcher with underline or pill styling"},
	"ThemeSwitch":    {Name: "ThemeSwitch", Category: "Navigation", Description: "Light and dark mode theme toggle switch"},

	// Feedback & Alerts
	"Alert":   {Name: "Alert", Category: "Feedback", Description: "Contextual alert box with info, success, warning, and error types"},
	"Callout": {Name: "Callout", Category: "Feedback", Description: "Prominent callout banner for highlights, tips, or quotes"},

	// Overlays
	"Collapse": {Name: "Collapse", Category: "Overlays", Description: "Expandable accordion disclosure block with smooth animation"},
	"Drawer":   {Name: "Drawer", Category: "Overlays", Description: "Slide-out panel overlay anchored to left or right screen edge"},
	"Modal":    {Name: "Modal", Category: "Overlays", Description: "Accessible dialog modal overlay with header, body, and footer slots"},
	"Tooltip":  {Name: "Tooltip", Category: "Overlays", Description: "Contextual hover tooltip bubble"},

	// Data Display & Complex
	"Card":       {Name: "Card", Category: "Data Display", Description: "Content card container with header, media, body, and footer slots"},
	"CodeWindow": {Name: "CodeWindow", Category: "Data Display", Description: "Syntax-highlighted code display window with copy button and titlebar"},
	"Ecosystem":  {Name: "Ecosystem", Category: "Data Display", Description: "Interactive ecosystem showcase grid with category filtering"},
	"Footer":     {Name: "Footer", Category: "Data Display", Description: "Full-width site footer with customizable link columns and copyright"},
	"List":       {Name: "List", Category: "Data Display", Description: "Styled bulleted or numbered list container with custom markers"},
	"Table":      {Name: "Table", Category: "Data Display", Description: "Responsive data table with striped, bordered, and hover states"},
	"Terminal":   {Name: "Terminal", Category: "Data Display", Description: "Retro-modern CLI terminal window display with animated prompts"},

	// Document
	"Head": {Name: "Head", Category: "Document", Description: "HTML head document tags injector for titles, meta, and styles"},
}

// List returns a sorted slice of all available component names.
func List() []string {
	names := make([]string, 0, len(Registry))
	for name := range Registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Get returns the raw .kiw template content for the named component.
// Name can be provided as "Button" or "Button.kiw".
func Get(name string) ([]byte, error) {
	name = strings.TrimSuffix(name, ".kiw")
	filename := name + ".kiw"
	body, err := FS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("component %q not found in Krewire UI: %w", name, err)
	}
	return body, nil
}
