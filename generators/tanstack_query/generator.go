package tanstackquery

import (
	"embed"
	"strings"

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
			"dependencies": map[string]interface{}{
				"@tanstack/react-query":          "^5.80.7",
				"@tanstack/react-query-devtools": "^5.80.7",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(fs, nil); err != nil {
		return err
	}

	// Inject QueryProvider into main.tsx.
	// Works regardless of which router (none, react-router, tanstack-router) ran first.
	if f, ok := ctx.State.GetFile("src/main.tsx"); ok {
		content := string(f.Content)
		if !strings.Contains(content, "QueryProvider") {
			content = strings.Replace(content,
				`import ReactDOM from "react-dom/client";`,
				"import ReactDOM from \"react-dom/client\";\nimport { QueryProvider } from \"./providers/query-provider\";",
				1)
			content = strings.Replace(content,
				"  <React.StrictMode>\n    ",
				"  <React.StrictMode>\n    <QueryProvider>\n      ",
				1)
			content = strings.Replace(content,
				"\n  </React.StrictMode>",
				"\n    </QueryProvider>\n  </React.StrictMode>",
				1)
			ctx.State.WriteFile("src/main.tsx", []byte(content), state.ContentRaw)
		}
	}

	return nil
}
