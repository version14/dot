package drizzletypescriptdeps

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
				"db:generate": "drizzle-kit generate",
				"db:migrate":  "drizzle-kit migrate",
				"db:push":     "drizzle-kit push",
				"db:studio":   "drizzle-kit studio",
			},
			"dependencies":    deps.NPM("drizzle-orm"),
			"devDependencies": deps.NPM("drizzle-kit"),
		})
		return nil
	})
}
