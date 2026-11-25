# `editors` (`github.com/krewire/ui/editors`)

## Purpose & Design Role

`editors` provides syntax highlighting, bracket matching, and folding definitions for `.kiw` single-file components. It ships a standard TextMate grammar (`scopeName: "text.html.kiw"`) packaged as a VS Code extension under `editors/vscode/`, along with a canonical sample file (`vscode/samples/Sample.kiw`) and automated Go tests that keep the grammar valid and synchronized with `kiw.Lint`.

Because TextMate grammars are editor-agnostic JSON, `kiw.tmLanguage.json` can also be loaded directly by Sublime Text, Zed, Helix, IntelliJ, and GitHub Linguist-style highlighters.

## Architectural Position

- **Layer**: Layer 2 (Developer Experience & Tooling Support within `github.com/krewire/ui`).
- **Distribution Model**: Model C Go module + static TextMate/VS Code extension assets.
- **Dependencies**: Go standard library (`encoding/json`, `os`, `regexp`, `strings`, `testing`) and `github.com/krewire/ui/kiw` (used in `editors_test.go` to validate `Sample.kiw`).

## Core Concepts & Assets

### VS Code Extension & TextMate Grammar (`vscode/`)
- `vscode/package.json`: Declares the VS Code extension metadata and registers the `text.html.kiw` grammar contribution.
- `vscode/syntaxes/kiw.tmLanguage.json`: TextMate grammar covering:
  - `frontmatter`: YAML frontmatter between leading `---` fences.
  - `blocks`: Template control flow (`{#if}`, `{:else if}`, `{:else}`, `{/if}`, `{#each items as item, i}`, `{/each}`).
  - `interpolation` & `template-action`: Single-brace `{prop}` interpolation and Go `{{ ... }}` actions.
  - `component-tag`: PascalCase component tags (`<Card ...>`) with props, `{expr}` values, and directives (`client:load`, `class:name={cond}`, `slot`).
  - `style`, `script`, `markdown`: Embedded CSS, JS/TS (`lang="ts"`), and Markdown blocks.
- `vscode/language-configuration.json`: Configures bracket matching, auto-closing pairs, and folding markers for `{#if}...{/if}` and `{#each}...{/each}`.
- `vscode/samples/Sample.kiw`: Reference `.kiw` file exercising every grammar construct.

### Automated Grammar Verification (`editors_test.go`)
- `TestGrammar(t *testing.T)`: Parses `kiw.tmLanguage.json`, checks `scopeName == "text.html.kiw"`, compiles every `match`, `begin`, and `end` regex pattern, verifies required repository rules exist, and validates `package.json` and `language-configuration.json`.
- `TestSampleLintsClean(t *testing.T)`: Runs `kiw.Lint` and `kiw.ParseKiw` against `vscode/samples/Sample.kiw` to ensure the sample remains valid `.kiw` syntax.

## Realistic Usage Examples

Installing the extension into VS Code via symlink for local development:

```bash
ln -s /home/krewire/workspace/krewire/ui/editors/vscode ~/.vscode/extensions/krewire.kiw-0.1.0
```

Or packaging a `.vsix` archive:

```bash
cd /home/krewire/workspace/krewire/ui/editors/vscode
npx @vscode/vsce package
code --install-extension kiw-0.1.0.vsix
```

## Verification

Run from `/home/krewire/workspace/krewire/ui/editors`:

```bash
gofmt -l .
go vet ./...
go test -v ./...
CGO_ENABLED=0 go build ./...
```
