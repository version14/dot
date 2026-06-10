package featureflagsposthog

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
				"posthog-js": "^1.384.2",
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
		posthog, _ := nextFS.ReadFile("next/posthog.ts")
		ctx.State.WriteFile("src/lib/posthog.ts", posthog, state.ContentRaw)
		env, _ := nextFS.ReadFile("next/.env.example")
		ctx.State.AppendFile(".env.example", env)
	} else {
		ctx.State.AppendFile(".env.example", []byte("# PostHog\nVITE_POSTHOG_KEY=phc_your_project_api_key\nVITE_POSTHOG_HOST=https://app.posthog.com\n"))
	}

	return nil
}
