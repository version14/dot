package prettierconfig

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile(".prettierrc", []byte(`{
  "semi": true,
  "singleQuote": true,
  "trailingComma": "all",
  "tabWidth": 2,
  "useTabs": false
}
`), state.ContentRaw)

	ctx.State.WriteFile(".prettierignore", []byte(`node_modules
dist
build
.env
`), state.ContentRaw)

	return nil
}
