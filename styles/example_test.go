package styles_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/krewire/ui/styles"
)

// Serve the canonical stylesheet from a product's HTTP server so pages can
// include it with <link rel="stylesheet" href="/styles/krewire-ui.css">.
func ExampleCSS() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /styles/krewire-ui.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		fmt.Fprint(w, styles.CSS)
	})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/styles/krewire-ui.css", nil))

	fmt.Println(rec.Code)
	fmt.Println(rec.Header().Get("Content-Type"))
	fmt.Println(strings.Contains(rec.Body.String(), "--pop-shadow"))
	// Output:
	// 200
	// text/css; charset=utf-8
	// true
}
