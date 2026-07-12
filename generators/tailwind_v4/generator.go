package tailwindv4

import (
	"embed"

	"github.com/version14/dot/internal/deps"
	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

//go:embed extra/vite.config.ts
var viteConfigBytes []byte

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies":    deps.NPM("tailwindcss"),
			"devDependencies": deps.NPM("@tailwindcss/vite", "@tailwindcss/postcss"),
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(filesFS, nil); err != nil {
		return err
	}

	if framework != "next" {
		ctx.State.WriteFile("vite.config.ts", viteConfigBytes, state.ContentRaw)
	}

	return nil
}
