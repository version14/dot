package prettierconfig

import (
	"strings"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

// defaultIgnoreEntries are the standard prettier ignore entries every project
// should have. Earlier generators (e.g. panda_css) may have pre-written entries
// to .prettierignore; we read them and merge so nothing is lost.
var defaultIgnoreEntries = []string{
	"node_modules",
	".dot",
	".next",
	"coverage",
	"dist",
	"build",
	"playwright-report",
	"storybook-static",
	"test-results",
	".env",
}

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile(".prettierrc", []byte(`{
  "semi": true,
  "singleQuote": true,
  "trailingComma": "all",
  "tabWidth": 2,
  "useTabs": false
}
`), state.ContentRaw)

	// Merge existing .prettierignore entries (written by earlier generators)
	// with the standard defaults, deduplicating. Order: existing first, then defaults.
	existing := ""
	if f, ok := ctx.State.GetFile(".prettierignore"); ok {
		existing = string(f.Content)
	}
	seen := map[string]bool{}
	var result []string
	for _, line := range strings.Split(existing, "\n") {
		line = strings.TrimRight(line, "\r")
		if line != "" && !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	for _, entry := range defaultIgnoreEntries {
		if !seen[entry] {
			seen[entry] = true
			result = append(result, entry)
		}
	}
	ctx.State.WriteFile(".prettierignore", []byte(strings.Join(result, "\n")+"\n"), state.ContentRaw)

	return nil
}
