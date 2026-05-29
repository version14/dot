package analyticsplausible

import (
	"embed"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

//go:embed all:next
var nextFS embed.FS

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"plausible-tracker": "^0.3.9",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(filesFS, nil); err != nil {
		return err
	}

	if framework == "next" {
		plausible, _ := nextFS.ReadFile("next/plausible.ts")
		ctx.State.WriteFile("src/lib/plausible.ts", plausible, state.ContentRaw)
		env, _ := nextFS.ReadFile("next/.env.example")
		ctx.State.AppendFile(".env.example", env)
	} else {
		ctx.State.AppendFile(".env.example", []byte("# Plausible Analytics\nVITE_PLAUSIBLE_DOMAIN=yourdomain.com\nVITE_PLAUSIBLE_API_HOST=https://plausible.io\n"))
	}

	return nil
}
