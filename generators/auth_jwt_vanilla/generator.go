package authjwtvanilla

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

//go:embed all:files
var fs embed.FS

func (g *Generator) Generate(ctx *dotapi.Context) error {
	renderer := render.NewLocalFolderRenderer(ctx.State)
	if err := renderer.Render(fs, nil); err != nil {
		return err
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"jsonwebtoken": "^9.0.2",
			},
			"devDependencies": map[string]interface{}{
				"@types/jsonwebtoken": "^9.0.7",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	// Append JWT env vars to .env.example
	existing := ""
	if f, ok := ctx.State.GetFile(".env.example"); ok {
		existing = string(f.Content)
	}
	updated := existing + fmt.Sprintf("\n# Auth (JWT)\nJWT_SECRET=%s\nJWT_EXPIRES_IN=7d\nJWT_REFRESH_EXPIRES_IN=30d\n", "change-me-to-a-random-secret")
	ctx.State.WriteFile(".env.example", []byte(updated), state.ContentRaw)

	return nil
}
