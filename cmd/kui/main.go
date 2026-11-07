package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/krewire/ui/components"
	"github.com/krewire/ui/styles"
)

const version = "0.2.0"

type Config struct {
	Schema        string `json:"$schema,omitempty"`
	Style         string `json:"style"`
	ComponentsDir string `json:"componentsDir"`
	CSSDir        string `json:"cssDir"`
}

func defaultConfig() Config {
	return Config{
		Schema:        "https://ui.krewire.com/schema.json",
		Style:         "default",
		ComponentsDir: "components",
		CSSDir:        "static/css",
	}
}

func loadConfig() Config {
	cfg := defaultConfig()
	data, err := os.ReadFile("kui.json")
	if err == nil {
		_ = json.Unmarshal(data, &cfg)
	}
	return cfg
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "init":
		runInit(args)
	case "list", "ls":
		runList(args)
	case "add":
		runAdd(args)
	case "info":
		runInfo(args)
	case "version", "-v", "--version":
		fmt.Printf("kui v%s (Krewire UI — Zero-dependency .kiw design system)\n", version)
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nRun 'kui help' for usage.\n", cmd)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`kui v%s — The Krewire UI CLI

Like shadcn/ui, but natively built for the .kiw DSL.
Own your components, copy-paste what you need, customize everything.

Usage:
  kui <command> [arguments]

Commands:
  init             Initialize Krewire UI in your project (creates kui.json & css)
  list             List all available .kiw components by category
  add <names...>   Copy specified .kiw components into your project
  info <name>      Inspect component details and template source
  version          Print kui version
  help             Show this help message

Examples:
  kui init
  kui list
  kui add Button Navbar HStack VStack
  kui add --all
  kui info Card
`, version)
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	componentsDir := fs.String("dir", "components", "Directory where .kiw components will be saved")
	cssDir := fs.String("css", "static/css", "Directory where krewire-ui.css will be saved")
	force := fs.Bool("force", false, "Overwrite existing kui.json and CSS")
	_ = fs.Parse(args)

	cfg := Config{
		Schema:        "https://ui.krewire.com/schema.json",
		Style:         "default",
		ComponentsDir: *componentsDir,
		CSSDir:        *cssDir,
	}

	// 1. Write kui.json
	if _, err := os.Stat("kui.json"); err == nil && !*force {
		fmt.Println("ℹ️  kui.json already exists. Use --force to overwrite.")
	} else {
		data, _ := json.MarshalIndent(cfg, "", "  ")
		if err := os.WriteFile("kui.json", data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing kui.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Created kui.json")
	}

	// 2. Ensure components directory exists
	if err := os.MkdirAll(cfg.ComponentsDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating components directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created components directory: %s/\n", cfg.ComponentsDir)

	// 3. Write krewire-ui.css
	if err := os.MkdirAll(cfg.CSSDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating CSS directory: %v\n", err)
		os.Exit(1)
	}
	cssTarget := filepath.Join(cfg.CSSDir, "krewire-ui.css")
	if _, err := os.Stat(cssTarget); err == nil && !*force {
		fmt.Printf("ℹ️  %s already exists. Use --force to overwrite.\n", cssTarget)
	} else {
		if err := os.WriteFile(cssTarget, []byte(styles.CSS), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", cssTarget, err)
			os.Exit(1)
		}
		fmt.Printf("✓ Installed stylesheet: %s\n", cssTarget)
	}

	fmt.Println("\n🎉 Krewire UI initialized successfully!")
	fmt.Println("Next steps:")
	fmt.Printf("  • Add components: kui add Button Card Navbar\n")
	fmt.Printf("  • Reference stylesheet in your layouts: <link rel=\"stylesheet\" href=\"/%s/krewire-ui.css\">\n", cfg.CSSDir)
}

func runList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	_ = fs.Parse(args)

	if *jsonOutput {
		all := make([]components.ComponentMeta, 0, len(components.Registry))
		for _, m := range components.Registry {
			all = append(all, m)
		}
		sort.Slice(all, func(i, j int) bool { return all[i].Name < all[j].Name })
		data, _ := json.MarshalIndent(all, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Group by category
	byCategory := make(map[string][]components.ComponentMeta)
	for _, m := range components.Registry {
		byCategory[m.Category] = append(byCategory[m.Category], m)
	}

	categories := make([]string, 0, len(byCategory))
	for cat := range byCategory {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	fmt.Printf("Available Krewire UI Components (%d total):\n\n", len(components.Registry))

	for _, cat := range categories {
		items := byCategory[cat]
		sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
		fmt.Printf("■ %s (%d)\n", cat, len(items))
		for _, item := range items {
			fmt.Printf("  • %-16s %s\n", item.Name, item.Description)
		}
		fmt.Println()
	}
	fmt.Println("Install components with: kui add <name...>")
}

func runAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	allFlag := fs.Bool("all", false, "Install all 51 components")
	overwrite := fs.Bool("overwrite", false, "Overwrite existing files")
	targetDir := fs.String("dir", "", "Target directory (defaults to kui.json componentsDir or 'components')")
	_ = fs.Parse(args)

	cfg := loadConfig()
	dir := cfg.ComponentsDir
	if *targetDir != "" {
		dir = *targetDir
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
		os.Exit(1)
	}

	var toAdd []string
	if *allFlag {
		toAdd = components.List()
	} else {
		toAdd = fs.Args()
	}

	if len(toAdd) == 0 {
		fmt.Println("Please specify component names to add, e.g.:")
		fmt.Println("  kui add Button Navbar HStack VStack")
		fmt.Println("Or install all components with: kui add --all")
		os.Exit(1)
	}

	var addedCount, skippedCount int

	for _, reqName := range toAdd {
		cleanName := resolveComponentName(reqName)
		if cleanName == "" {
			fmt.Fprintf(os.Stderr, "⚠️  Component %q not found in registry (run 'kui list' to view available)\n", reqName)
			skippedCount++
			continue
		}

		targetPath := filepath.Join(dir, cleanName+".kiw")
		if _, err := os.Stat(targetPath); err == nil && !*overwrite {
			fmt.Printf("ℹ️  %s already exists. Use --overwrite to replace.\n", targetPath)
			skippedCount++
			continue
		}

		body, err := components.Get(cleanName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching %s: %v\n", cleanName, err)
			skippedCount++
			continue
		}

		if err := os.WriteFile(targetPath, body, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", targetPath, err)
			skippedCount++
			continue
		}

		fmt.Printf("✓ Added %s -> %s\n", cleanName, targetPath)
		addedCount++
	}

	fmt.Printf("\nDone: %d component(s) added, %d skipped.\n", addedCount, skippedCount)
}

func runInfo(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: kui info <component-name>")
		os.Exit(1)
	}

	name := resolveComponentName(args[0])
	if name == "" {
		fmt.Fprintf(os.Stderr, "Component %q not found in registry.\n", args[0])
		os.Exit(1)
	}

	meta := components.Registry[name]
	body, err := components.Get(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading component %s: %v\n", name, err)
		os.Exit(1)
	}

	fmt.Printf("Component:   %s\n", meta.Name)
	fmt.Printf("Category:    %s\n", meta.Category)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Printf("File size:   %d bytes\n", len(body))
	fmt.Println("\nSource Preview (.kiw):")
	fmt.Println("----------------------------------------")
	lines := strings.Split(string(body), "\n")
	maxPreview := 30
	if len(lines) < maxPreview {
		maxPreview = len(lines)
	}
	for i := 0; i < maxPreview; i++ {
		fmt.Println(lines[i])
	}
	if len(lines) > maxPreview {
		fmt.Printf("... (%d more lines)\n", len(lines)-maxPreview)
	}
	fmt.Println("----------------------------------------")
}

func resolveComponentName(input string) string {
	input = strings.TrimSuffix(input, ".kiw")
	for name := range components.Registry {
		if strings.EqualFold(name, input) {
			return name
		}
	}
	return ""
}
