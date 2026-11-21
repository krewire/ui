package kiw

import (
	"reflect"
	"testing"
)

func TestPropAliases(t *testing.T) {
	cases := map[string][]string{
		"variant":       {"variant", "Variant"},
		"Variant":       {"Variant", "variant"},
		"aria-label":    {"aria-label", "AriaLabel", "ariaLabel"},
		"show_value":    {"show_value", "ShowValue", "showValue"},
		"property-name": {"property-name", "PropertyName", "propertyName"},
		"brand-name":    {"brand-name", "BrandName", "brandName"},
		"id":            {"id", "Id", "ID"},
		"":              nil,
	}
	for in, want := range cases {
		if got := PropAliases(in); !reflect.DeepEqual(got, want) {
			t.Errorf("PropAliases(%q) = %v, want %v", in, got, want)
		}
	}

	// Verify kebab-case with initialisms / github
	ghAliases := PropAliases("github-url")
	mustContain := []string{"github-url", "GithubUrl", "githubUrl", "gitHubURL", "githubURL", "GitHubURL"}
	for _, expected := range mustContain {
		found := false
		for _, a := range ghAliases {
			if a == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("PropAliases(\"github-url\") missing alias %q, got: %v", expected, ghAliases)
		}
	}
}
