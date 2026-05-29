package prettierfrontendrules

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	return ctx.State.UpdateJSON(".prettierrc", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"semi":            true,
			"singleQuote":     false,
			"jsxSingleQuote":  false,
			"trailingComma":   "all",
			"printWidth":      80,
			"tabWidth":        2,
			"bracketSpacing":  true,
			"bracketSameLine": false,
		})
		return nil
	})
}
