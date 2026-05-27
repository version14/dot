package cssmodules

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("src/styles/global.css", []byte(globalCSS), state.ContentRaw)
	ctx.State.WriteFile("src/styles/App.module.css", []byte(appModuleCSS), state.ContentRaw)
	return nil
}

const globalCSS = `*,
*::before,
*::after {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: system-ui, sans-serif;
}
`

const appModuleCSS = `.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.title {
  font-size: 2rem;
  font-weight: bold;
}
`
