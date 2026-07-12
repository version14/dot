package vitesttestinglibrary

import (
	"embed"

	"github.com/version14/dot/internal/deps"
	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var fs embed.FS

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"test":          "vitest run",
				"test:watch":    "vitest",
				"test:coverage": "vitest run --coverage",
			},
			"devDependencies": deps.NPM("vitest", "@vitest/coverage-v8", "@testing-library/react", "@testing-library/user-event", "@testing-library/jest-dom", "jsdom"),
		})
		return nil
	}); err != nil {
		return err
	}

	return render.NewLocalFolderRenderer(ctx.State).Render(fs, nil)
}
