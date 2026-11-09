package components_test

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"testing/fstest"

	"github.com/krewire/ui/components"
)

// Read one component's .kiw source. The name may be given with or without the
// .kiw extension.
func ExampleGet() {
	src, err := components.Get("Button")
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.Contains(string(src), "<style>"))

	_, err = components.Get("NoSuchComponent")
	fmt.Println(errors.Is(err, fs.ErrNotExist))
	// Output:
	// true
	// true
}

// List returns every component name in sorted order.
func ExampleList() {
	names := components.List()
	fmt.Println(len(names))
	fmt.Println(names[:3])
	// Output:
	// 66
	// [Accordion ActionMenu Alert]
}

// Registry carries the metadata used to build catalogs and documentation.
func ExampleRegistry() {
	counts := map[string]int{}
	for _, meta := range components.Registry {
		counts[meta.Category]++
	}
	for _, c := range []string{"Layout", "Atoms", "Forms", "Navigation", "Feedback", "Overlays", "Data Display", "Document"} {
		fmt.Printf("%s: %d\n", c, counts[c])
	}
	// Output:
	// Layout: 8
	// Atoms: 12
	// Forms: 10
	// Navigation: 11
	// Feedback: 4
	// Overlays: 7
	// Data Display: 13
	// Document: 1
}

// FS is a standard io/fs filesystem, so the usual fs helpers work on it.
func ExampleFS() {
	stacks, err := fs.Glob(components.FS, "*Stack.kiw")
	if err != nil {
		panic(err)
	}
	fmt.Println(stacks)
	// Output:
	// [HStack.kiw Stack.kiw VStack.kiw]
}

// A product composes its own component set on top of kui: product components
// win over kui components with the same name, and everything else falls back to
// the shared catalog. This is the same precedence the .kiw site loader applies
// to a project's components/ directory.
func Example_productOverride() {
	product := fstest.MapFS{
		// A product-owned component that composes kui primitives.
		"ServerCard.kiw": {Data: []byte(`<Card title={name}><Badge variant="success">{status}</Badge></Card>`)},
		// A product-specific replacement for a kui component.
		"Brand.kiw": {Data: []byte(`<a href="/" class="brand">Acme</a>`)},
	}

	resolve := func(name string) (string, []byte) {
		if body, err := fs.ReadFile(product, name+".kiw"); err == nil {
			return "product", body
		}
		body, err := components.Get(name)
		if err != nil {
			return "missing", nil
		}
		return "kui", body
	}

	names := []string{"ServerCard", "Brand", "Card", "Badge", "Unknown"}
	sort.Strings(names)
	for _, n := range names {
		src, _ := resolve(n)
		fmt.Printf("%s: %s\n", n, src)
	}
	// Output:
	// Badge: kui
	// Brand: product
	// Card: kui
	// ServerCard: product
	// Unknown: missing
}
