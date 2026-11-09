package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/krewire/ui/components"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(NewHandler(NewCatalog(productComponents())))
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url string) (*http.Response, string) {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res, string(body)
}

func TestServesStylesheet(t *testing.T) {
	srv := newTestServer(t)
	res, body := get(t, srv.URL+"/styles/krewire-ui.css")
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/css") {
		t.Errorf("Content-Type = %q", ct)
	}
	if !strings.Contains(body, "--primary") {
		t.Error("stylesheet does not contain the --primary token")
	}
}

func TestResolvePrefersProductComponents(t *testing.T) {
	srv := newTestServer(t)
	cases := []struct {
		file   string
		source Source
	}{
		{"ServerCard.kiw", SourceProduct}, // product-owned
		{"Button.kiw", SourceKUI},         // shared foundation
	}
	for _, tc := range cases {
		res, body := get(t, srv.URL+"/components/"+tc.file)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("%s: status = %d", tc.file, res.StatusCode)
		}
		if got := Source(res.Header.Get("X-Component-Source")); got != tc.source {
			t.Errorf("%s: source = %q, want %q", tc.file, got, tc.source)
		}
		if strings.TrimSpace(body) == "" {
			t.Errorf("%s: empty body", tc.file)
		}
	}

	res, _ := get(t, srv.URL+"/components/DoesNotExist.kiw")
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("unknown component: status = %d, want 404", res.StatusCode)
	}
}

func TestCatalogListsBothLayers(t *testing.T) {
	srv := newTestServer(t)
	_, body := get(t, srv.URL+"/components")
	var list []struct {
		Name   string `json:"name"`
		Source Source `json:"source"`
	}
	if err := json.Unmarshal([]byte(body), &list); err != nil {
		t.Fatal(err)
	}
	if want := len(components.Registry) + 1; len(list) != want {
		t.Fatalf("catalog has %d entries, want %d (kui + ServerCard)", len(list), want)
	}
}

// componentTagRe matches a JSX-like component tag (<Name ...>, <Name/>).
var componentTagRe = regexp.MustCompile(`<([A-Z][A-Za-z0-9]*)[\s/>]`)

// TestSiteUsesOnlyKnownComponents checks that every component tag in the
// example site resolves to either a kui component or a product component, so
// the example cannot silently drift from the kui catalog.
func TestSiteUsesOnlyKnownComponents(t *testing.T) {
	known := NewCatalog(productComponents()).Names()
	err := fs.WalkDir(site, "site", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".kiw") {
			return err
		}
		src, err := fs.ReadFile(site, path)
		if err != nil {
			return err
		}
		for _, m := range componentTagRe.FindAllStringSubmatch(string(src), -1) {
			if _, ok := known[m[1]]; !ok {
				t.Errorf("%s: <%s> is neither a kui nor a product component", path, m[1])
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
