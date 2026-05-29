package analyticsga4

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
				"react-ga4": "^3.0.1",
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
		ga4, _ := nextFS.ReadFile("next/ga4.ts")
		ctx.State.WriteFile("src/lib/ga4.ts", ga4, state.ContentRaw)
		env, _ := nextFS.ReadFile("next/.env.example")
		ctx.State.AppendFile(".env.example", env)
	} else {
		ctx.State.AppendFile(".env.example", []byte("# Google Analytics 4\nVITE_GA_MEASUREMENT_ID=G-XXXXXXXXXX\n"))
	}

	return nil
}
