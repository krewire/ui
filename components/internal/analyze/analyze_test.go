package analyze

import (
	"reflect"
	"testing"
)

func TestPropName(t *testing.T) {
	for in, want := range map[string]string{
		"Variant": "variant", ".Variant": "variant", "ID": "id", "URLPath": "urlPath",
		"GitHubURL": "gitHubURL", "brandName": "brandName", "Item.Label": "item",
	} {
		if got := PropName(in); got != want {
			t.Errorf("PropName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAnalyze(t *testing.T) {
	src := `{#if href}<a href="{href}" class="{{class "btn" .Variant "primary" .Size "md" .Class}}">{body}</a>{/if}
{#each items as item, i}<Badge label="{item.label}" />{/each}
<style>.btn{color:red}</style>`
	f := Analyze("Button", src)
	if want := []string{"body", "class", "href", "items", "size", "variant"}; !reflect.DeepEqual(f.Props, want) {
		t.Errorf("Props = %v, want %v", f.Props, want)
	}
	if want := map[string]string{"variant": "primary", "size": "md"}; !reflect.DeepEqual(f.Defaults, want) {
		t.Errorf("Defaults = %v, want %v", f.Defaults, want)
	}
	if want := []string{"Badge"}; !reflect.DeepEqual(f.Dependencies, want) {
		t.Errorf("Dependencies = %v, want %v", f.Dependencies, want)
	}
	if !f.Styles || f.Script {
		t.Errorf("Styles/Script = %v/%v, want true/false", f.Styles, f.Script)
	}
}
