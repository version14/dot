package biomeconfig

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("biome.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"$schema": "https://biomejs.dev/schemas/2.4.16/schema.json",
			"linter": map[string]interface{}{
				"enabled": true,
				"rules":   map[string]interface{}{"recommended": true},
			},
			"formatter": map[string]interface{}{
				"enabled":     true,
				"indentStyle": "space",
				"indentWidth": 2,
			},
			"files": map[string]interface{}{
				"experimentalScannerIgnores": []interface{}{".dot/"},
			},
		})
		return nil
	}); err != nil {
		return err
	}

	return ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"lint":   "biome check .",
				"format": "biome format --write .",
			},
			"devDependencies": map[string]interface{}{
				"@biomejs/biome": "^2.4.16",
			},
		})
		return nil
	})
}
