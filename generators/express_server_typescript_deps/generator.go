package expressservertypescriptdeps

import (
	"github.com/version14/dot/internal/deps"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	return ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"dev":   "nodemon --exec tsx src/index.ts",
				"build": "tsc",
				"start": "node dist/index.js",
			},
			"dependencies":    deps.NPM("cors", "dotenv", "express"),
			"devDependencies": deps.NPM("@types/cors", "@types/express", "@types/node", "nodemon", "tsx"),
		})
		return nil
	})
}
