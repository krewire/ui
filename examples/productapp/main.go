// Command productapp is a minimal product application that consumes Krewire UI
// (kui) as a plain Go module.
//
// It shows the integration boundary described in ARCHITECTURE.md:
//
//   - The app imports only github.com/krewire/ui/components and
//     github.com/krewire/ui/styles. It does not import the kiw CLI or any other
//     part of github.com/krewire/krewire.
//   - Product-owned components (site/components) live in the product and take
//     precedence over kui components with the same name, mirroring the override
//     rule used by the .kiw site loader.
//   - Product screens (site/pages, site/layouts) compose kui primitives and
//     product components; kui itself knows nothing about servers or billing.
//
// The site/ directory is a regular .kiw site that a .kiw renderer (for example
// `kiw build` in github.com/krewire/krewire) can render. This program only
// serves the stylesheet and the resolved component sources, so it stays
// renderer-agnostic.
//
// Run it with:
//
//	go run ./examples/productapp -addr :8080
package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/krewire/ui/components"
	"github.com/krewire/ui/styles"
)

//go:embed site
var site embed.FS

// Source identifies where a resolved component came from.
type Source string

const (
	// SourceKUI marks a component served from the shared kui catalog.
	SourceKUI Source = "kui"
	// SourceProduct marks a component owned by the product application.
	SourceProduct Source = "product"
)

// Catalog resolves component names against product components first and the
// kui catalog second.
type Catalog struct {
	product fs.FS
}

// NewCatalog returns a Catalog whose product layer is the given filesystem of
// <Name>.kiw files. A nil product filesystem yields the plain kui catalog.
func NewCatalog(product fs.FS) *Catalog {
	return &Catalog{product: product}
}

// Resolve returns the .kiw source for name and where it was found.
func (c *Catalog) Resolve(name string) ([]byte, Source, error) {
	name = strings.TrimSuffix(name, ".kiw")
	if c.product != nil {
		body, err := fs.ReadFile(c.product, name+".kiw")
		if err == nil {
			return body, SourceProduct, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, "", err
		}
	}
	body, err := components.Get(name)
	if err != nil {
		return nil, "", err
	}
	return body, SourceKUI, nil
}

// Names maps every resolvable component name to its source.
func (c *Catalog) Names() map[string]Source {
	out := make(map[string]Source, len(components.Registry))
	for _, name := range components.List() {
		out[name] = SourceKUI
	}
	if c.product != nil {
		matches, _ := fs.Glob(c.product, "*.kiw")
		for _, m := range matches {
			out[strings.TrimSuffix(m, ".kiw")] = SourceProduct
		}
	}
	return out
}

// NewHandler serves the kui stylesheet and the resolved component catalog.
//
//	GET /styles/krewire-ui.css   the canonical kui stylesheet
//	GET /components              JSON list of {name, source}
//	GET /components/{name}.kiw   resolved component source
func NewHandler(c *Catalog) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /styles/krewire-ui.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		_, _ = w.Write([]byte(styles.CSS))
	})
	mux.HandleFunc("GET /components", func(w http.ResponseWriter, r *http.Request) {
		type entry struct {
			Name   string `json:"name"`
			Source Source `json:"source"`
		}
		names := c.Names()
		list := make([]entry, 0, len(names))
		for n, s := range names {
			list = append(list, entry{Name: n, Source: s})
		}
		sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	})
	mux.HandleFunc("GET /components/{file}", func(w http.ResponseWriter, r *http.Request) {
		file := r.PathValue("file")
		if !strings.HasSuffix(file, ".kiw") {
			http.NotFound(w, r)
			return
		}
		body, src, err := c.Resolve(file)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Component-Source", string(src))
		_, _ = w.Write(body)
	})
	return mux
}

// productComponents returns the product's own component directory.
func productComponents() fs.FS {
	sub, err := fs.Sub(site, "site/components")
	if err != nil {
		panic(err) // the directory is embedded at build time
	}
	return sub
}

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	log.Printf("productapp listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, NewHandler(NewCatalog(productComponents()))))
}
