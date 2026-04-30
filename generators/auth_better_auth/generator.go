package authbetterauth

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
				"better-auth": "^1.2.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	// Append BETTER_AUTH_SECRET to .env.example
	existing := ""
	if f, ok := ctx.State.GetFile(".env.example"); ok {
		existing = string(f.Content)
	}
	updated := existing + fmt.Sprintf("\n# Auth (BetterAuth)\nBETTER_AUTH_SECRET=%s\nBETTER_AUTH_URL=http://localhost:${PORT:-3000}\n", generateSecretPlaceholder())
	ctx.State.WriteFile(".env.example", []byte(updated), state.ContentRaw)

	return nil
}

func generateSecretPlaceholder() string {
	return "change-me-to-a-random-32-char-secret"
}
