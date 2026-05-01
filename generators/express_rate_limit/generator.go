package expressratelimit

import (
	"embed"
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

const limiterImport = "import { limiter } from './shared/middlewares/rate-limit.middleware';\n"
const limiterUse = "app.use(limiter);\n"

func (g *Generator) Generate(ctx *dotapi.Context) error {
	renderer := render.NewLocalFolderRenderer(ctx.State)
	if err := renderer.Render(fs, nil); err != nil {
		return err
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"express-rate-limit": "^7.4.1",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	// Inject limiter into app.ts after app.use(express.json())
	f, ok := ctx.State.GetFile("src/app.ts")
	if !ok {
		return nil
	}
	content := string(f.Content)
	if strings.Contains(content, "limiter") {
		return nil
	}

	content = limiterImport + content
	// Apply limiter after cors/json middleware setup, before routes
	content = strings.Replace(content, "app.use(express.urlencoded", limiterUse+"\napp.use(express.urlencoded", 1)

	ctx.State.WriteFile("src/app.ts", []byte(content), state.ContentRaw)
	return nil
}
