package authbetterauth

import (
	"embed"
	"fmt"
	"strings"

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
				"better-auth":   "^1.2.0",
				"cookie-parser": "^1.4.7",
			},
			"devDependencies": map[string]interface{}{
				"@types/cookie-parser": "^1.4.8",
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

	// Inject cookie-parser middleware into app.ts
	if f, ok := ctx.State.GetFile("src/app.ts"); ok {
		content := string(f.Content)
		if !strings.Contains(content, "cookieParser") {
			content = "import cookieParser from 'cookie-parser';\n" + content
			content = strings.Replace(content, "app.use(express.json());", "app.use(express.json());\napp.use(cookieParser());", 1)
			ctx.State.WriteFile("src/app.ts", []byte(content), state.ContentRaw)
		}
	}

	return nil
}

func generateSecretPlaceholder() string {
	return "change-me-to-a-random-32-char-secret"
}
